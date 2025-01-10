// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hcx "github.com/vmware/terraform-provider-hcx/hcx"
)

const (
	vmcMaxRetries               = 12
	vmcRetryInterval            = 10 * time.Second
	vmcMaxRetryInterval         = 5 * time.Minute
	vmcActivationStatus         = "ACTIVE"
	vmcActivationFailedStatus   = "ACTIVATION_FAILED"
	vmcDeactivationStatus       = "DE-ACTIVATED"
	vmcDeactivationFailedStatus = "DEACTIVATION_FAILED"
)

// resourceVmc defines the VMware Cloud SDDC resource for Terraform management, supporting creation, reading, updates,
// and deletion. It allows activation and deactivation of HCX on an SDDC, with fields to specify and retrieve SDDC
// details and related metadata.
func resourceVmc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVmcCreate,
		ReadContext:   resourceVmcRead,
		UpdateContext: resourceVmcUpdate,
		DeleteContext: resourceVmcDelete,

		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the SDDC. Either 'sddc_id' or 'sddc_name' must be specified.",
			},
			"sddc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SDDC. Either 'sddc_name' or 'sddc_id' must be specified.",
			},
			"cloud_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of HCX Cloud. Used for the site pairing configuration.",
			},
			"cloud_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the HCX Cloud.",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the HCX Cloud.",
			},
		},
	}
}

// resourceVmcCreate triggers the activation of HCX on a specified SDDC and ensures the process completes successfully.
// It authenticates, validates, and polls the activation status until the SDDC reaches the desired state.
// In case of repeated failures or errors during activation, it returns diagnostic information for troubleshooting.
func resourceVmcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcx.Client)

	// Authenticate and fetch the SDDC details.
	sddc, diags := authenticateAndFetchSDDC(ctx, d, client)
	if diags != nil {
		return diags
	}

	// Validate if the SDDC is already activated.
	if sddc.DeploymentStatus == vmcActivationStatus {
		return diag.Errorf("sddc is already activated")
	}

	// Trigger the activation of HCX on the specified SDDC.
	_, err := hcx.ActivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.Errorf("failed to activate hcx on the sddc: %s", err.Error())
	}

	// Poll the activation status until the process completes successfully or fails.
	for retries := 0; retries < vmcMaxRetries; retries++ {
		// Fetch the updated SDDC details.
		sddc, diags = authenticateAndFetchSDDC(ctx, d, client)
		if diags != nil {
			hclog.Default().Info("[INFO] - resourceVmcCreate() - error retrieving sddc status; retrying...",
				"sddc_id", sddc.ID, "sddc_name", sddc.Name, "error", diags[0].Summary, "RetryCount", retries)

			// Exit if the retry limit is exceeded.
			if retries == vmcMaxRetries-1 {
				return diags
			}

			// Apply exponential backoff with capped interval.
			waitTime := calculateBackoff(retries)
			time.Sleep(waitTime)
			continue
		}

		// Check the current deployment status of the SDDC.
		switch sddc.DeploymentStatus {
		case vmcActivationStatus:
			// Activation successful. Refresh the resource state by calling "read".
			return resourceVmcRead(ctx, d, m)

		case vmcActivationFailedStatus:
			// Explicit activation failure. Return an appropriate error message.
			return diag.Errorf("activation failed with status: %s", vmcActivationFailedStatus)

		default:
			hclog.Default().Warn("unknown sddc status during polling", "status", sddc.DeploymentStatus)
		}

		// Wait before the next polling attempt.
		time.Sleep(vmcRetryInterval)
	}

	// If retry limit is reached without resolution, return an error.
	return diag.Errorf("maximum retries reached while activating the sddc")
}

// resourceVmcRead reads the current state of the SDDC from the remote API and updates resource data.
// Returns diagnostic information on any errors encountered during authentication or data retrieval.
func resourceVmcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcx.Client)

	// Authenticate and fetch the SDDC details.
	sddc, diags := authenticateAndFetchSDDC(ctx, d, client)
	if diags != nil {
		return diags
	}

	// Set resource data fields with retrieved SDDC details.
	if err := d.Set("cloud_url", sddc.CloudURL); err != nil {
		return diag.Errorf("failed to set 'cloud_url': %s", err.Error())
	}
	if err := d.Set("cloud_name", sddc.CloudName); err != nil {
		return diag.Errorf("failed to set 'cloud_name': %s", err.Error())
	}
	if err := d.Set("cloud_type", sddc.CloudType); err != nil {
		return diag.Errorf("failed to set 'cloud_type': %s", err.Error())
	}
	d.SetId(sddc.ID)

	// Return diagnostics response.
	return nil
}

// resourceVmcUpdate handles the update logic for the resource, ensuring immutable fields are not modified.
// It validates updates to "sddc_id" and "sddc_name" and returns an error if modifications are attempted.
func resourceVmcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Validate immutable fields to prevent unsupported updates.
	if d.HasChange("sddc_id") || d.HasChange("sddc_name") {
		return diag.Errorf("opdating 'sddc_id' or 'sddc_name' is not supported")
	}

	// Return the current state using Read.
	return resourceVmcRead(ctx, d, m)
}

// resourceVmcDelete handles the deletion of an SDDC by deactivating HCX and polling status until completion or failure.
func resourceVmcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcx.Client)

	// Authenticate and fetch the SDDC details.
	sddc, diags := authenticateAndFetchSDDC(ctx, d, client)
	if diags != nil {
		return diags
	}

	// Trigger the deactivation of HCX on the specified SDDC.
	_, err := hcx.DeactivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.Errorf("failed to deactivate hcx on the sddc: %s", err.Error())
	}

	// Poll the deactivation status until the process completes successfully or fails.
	for retries := 0; retries < vmcMaxRetries; retries++ {
		// Fetch the updated status of the SDDC.
		sddc, diags = authenticateAndFetchSDDC(ctx, d, client)
		if diags != nil {
			hclog.Default().Info("[INFO] - resourceVmcDelete() - error retrieving sddc status; retrying...",
				"sddc_id", sddc.ID, "sddc_name", sddc.Name, "error", diags[0].Summary, "RetryCount", retries)

			// Exit if the retry limit is exceeded
			if retries == vmcMaxRetries-1 {
				return diags
			}

			// Apply exponential backoff with capped interval.
			waitTime := calculateBackoff(retries)
			time.Sleep(waitTime)
			continue
		}

		// Evaluate the current deployment status of the SDDC.
		switch sddc.DeploymentStatus {
		case vmcDeactivationStatus:
			// Deactivation successful. Exit the deletion process.
			return nil

		case vmcDeactivationFailedStatus:
			// Explicit deactivation failure. Return an appropriate error message.
			return diag.Errorf("deactivation failed with status: %s", vmcDeactivationFailedStatus)

		case "":
			// SDDC resource already missing or deleted.
			hclog.Default().Info("[INFO] - resourceVmcDelete() - sddc already missing or deleted.")
			return nil

		default:
			hclog.Default().Warn("unknown sddc status during polling", "status", sddc.DeploymentStatus)
		}
	}

	// If retry limit is reached without resolution, return an error.
	return diag.Errorf("maximum retries reached while deleting the sddc")
}

// authenticateAndFetchSDDC authenticates with VMware Cloud and HCX APIs, and retrieves SDDC details by name or ID.
func authenticateAndFetchSDDC(ctx context.Context, d *schema.ResourceData, client *hcx.Client) (hcx.SDDC, diag.Diagnostics) {
	// Extract input parameters.
	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	// Validate inputs.
	if sddcName == "" && sddcID == "" {
		return hcx.SDDC{}, diag.Errorf("neither 'sddc_name' nor 'sddc_id' was specified")
	}

	// Authenticate with VMware Cloud Services.
	accessToken, err := hcx.VmcAuthenticate(token)
	if err != nil {
		return hcx.SDDC{}, diag.Errorf("failed to authenticate with VMware Cloud Services: %s", err.Error())
	}

	// Authenticate the HCX API.
	err = hcx.HcxCloudAuthenticate(client, accessToken)
	if err != nil {
		return hcx.SDDC{}, diag.Errorf("failed to authenticate hcx api: %s", err.Error())
	}

	// Retrieve the SDDC details.
	var sddc hcx.SDDC
	if sddcID != "" {
		sddc, err = hcx.GetSddcByID(client, sddcID)
	} else {
		sddc, err = hcx.GetSddcByName(client, sddcName)
	}
	if err != nil {
		return hcx.SDDC{}, diag.Errorf("failed to retrieve sddc details: %s", err.Error())
	}

	return sddc, nil
}

// calculateBackoff computes the exponential backoff duration with a capped interval based on the number of retries.
func calculateBackoff(retries int) time.Duration {
	wait := vmcRetryInterval * time.Duration(1<<retries)
	if wait > vmcMaxRetryInterval {
		wait = vmcMaxRetryInterval
	}
	return wait
}
