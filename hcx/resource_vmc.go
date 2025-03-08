// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"math/rand"
	"time"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceVmc defines the resource schema for managing a VMware Cloud on AWS configuration.
func resourceVmc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVmcCreate,
		ReadContext:   resourceVmcRead,
		UpdateContext: resourceVmcUpdate,
		DeleteContext: resourceVmcDelete,

		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:         schema.TypeString,
				Description:  "The ID of the SDDC.",
				Optional:     true,
				ExactlyOneOf: []string{"sddc_id", "sddc_name"}, // Enforces that at least one of them is provided.
			},
			"sddc_name": {
				Type:         schema.TypeString,
				Description:  "The name of the SDDC.",
				Optional:     true,
				ExactlyOneOf: []string{"sddc_id", "sddc_name"}, // Enforces that at least one of them is provided.
			},
			"cloud_url": {
				Type:        schema.TypeString,
				Description: "The URL of HCX Cloud, used for the site pairing configuration.",
				Computed:    true,
			},
			"cloud_name": {
				Type:        schema.TypeString,
				Description: "The name of the HCX Cloud.",
				Computed:    true,
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Description: "The type of the HCX Cloud. Use 'nsp' for VMware Cloud on AWS.",
				Computed:    true,
			},
		},
	}
}

// resourceVmcCreate creates the VMware Cloud on AWS configuration.
func resourceVmcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Authenticate and fetch the SDDC details.
	sddc, diags := authenticateAndFetchSDDC(ctx, d, client)
	if diags != nil {
		return diags
	}

	// Validate if the SDDC is already activated.
	if sddc.DeploymentStatus == constants.VmcActivationActiveStatus {
		return diag.Errorf("sddc is already activated")
	}

	// Trigger the activation of HCX on the specified SDDC.
	_, err := ActivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.Errorf("failed to activate hcx on the sddc: %s", err.Error())
	}

	// Poll the activation status until the process completes successfully or fails.
	for retries := 0; retries < constants.VmcMaxRetries; retries++ {
		// Fetch the updated SDDC details.
		sddc, diags = authenticateAndFetchSDDC(ctx, d, client)
		if diags != nil {
			hclog.Default().Info("[INFO] - resourceVmcCreate() - error retrieving SDDC status; retrying...",
				"sddc_id", sddc.ID, "sddc_name", sddc.Name, "error", diags[0].Summary, "RetryCount", retries)

			// Exit if the retry limit is exceeded.
			if retries == constants.VmcMaxRetries-1 {
				return diags
			}

			// Backoff before retrying.
			waitTime := calculateBackoff(retries)
			if err := waitWithContext(ctx, waitTime); err != nil { // Respect context cancellations.
				return diag.FromErr(err)
			}
			continue
		}

		// Check the current deployment status of the SDDC.
		switch sddc.DeploymentStatus {
		case constants.VmcActivationActiveStatus:
			// Activation successful. Refresh the resource state by calling "read".
			return resourceVmcRead(ctx, d, m)

		case constants.VmcActivationFailedStatus:
			// Explicit activation failure. Return an appropriate error message.
			return diag.Errorf("activation failed with status: %s", constants.VmcActivationFailedStatus)

		default:
			hclog.Default().Warn("unknown SDDC status during polling", "status", sddc.DeploymentStatus)
		}

		// Wait before the next polling attempt.
		if err := waitWithContext(ctx, constants.VmcRetryInterval); err != nil {
			return diag.FromErr(err)
		}
	}

	// If retry limit is reached without resolution, return an error.
	return diag.Errorf("maximum retries reached while activating the SDDC")
}

// resourceVmcRead retrieves the VMware Cloud on AWS configuration.
func resourceVmcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

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

// resourceVmcUpdate updates the VMware Cloud on AWS resource configuration.
func resourceVmcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Validate immutable fields to prevent unsupported updates.
	if d.HasChange("sddc_id") || d.HasChange("sddc_name") {
		return diag.Errorf("opdating 'sddc_id' or 'sddc_name' is not supported")
	}

	// Return the current state using Read.
	return resourceVmcRead(ctx, d, m)
}

// resourceVmcDelete removes the VMware Cloud on AWS configuration and clears the state of the resource in the schema.
func resourceVmcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Authenticate and fetch the SDDC details.
	sddc, diags := authenticateAndFetchSDDC(ctx, d, client)
	if diags != nil {
		return diags
	}

	// Trigger the deactivation of HCX on the specified SDDC.
	_, err := DeactivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.Errorf("failed to deactivate hcx on the sddc: %s", err.Error())
	}

	// Poll the deactivation status until the process completes successfully or fails.
	for retries := 0; retries < constants.VmcMaxRetries; retries++ {
		// Fetch the updated status of the SDDC.
		sddc, diags = authenticateAndFetchSDDC(ctx, d, client)
		if diags != nil {
			hclog.Default().Info("[INFO] - resourceVmcDelete() - error retrieving SDDC status; retrying...",
				"sddc_id", sddc.ID, "sddc_name", sddc.Name, "error", diags[0].Summary, "RetryCount", retries)

			// Exit if the retry limit is exceeded.
			if retries == constants.VmcMaxRetries-1 {
				return diags
			}

			// Backoff before retrying.
			waitTime := calculateBackoff(retries)
			if err := waitWithContext(ctx, waitTime); err != nil {
				return diag.FromErr(err)
			}
			continue
		}

		// Evaluate the current deployment status of the SDDC.
		switch sddc.DeploymentStatus {
		case constants.VmcDeactivationInactiveStatus:
			// Deactivation successful. Exit the deletion process.
			return nil

		case constants.VmcDeactivationFailedStatus:
			// Explicit deactivation failure. Return an appropriate error message.
			return diag.Errorf("deactivation failed with status: %s", constants.VmcDeactivationFailedStatus)

		case "":
			// SDDC resource already missing or deleted.
			hclog.Default().Info("[INFO] - resourceVmcDelete() - SDDC already missing or deleted.")
			return nil

		default:
			hclog.Default().Warn("unknown SDDC status during polling", "status", sddc.DeploymentStatus)
		}

		// Wait before polling again.
		if err := waitWithContext(ctx, constants.VmcRetryInterval); err != nil {
			return diag.FromErr(err)
		}
	}

	// If retry limit is reached without resolution, return an error.
	return diag.Errorf("maximum retries reached while deleting the SDDC")
}

// authenticateAndFetchSDDC authenticates with VMware Cloud and HCX APIs, and retrieves SDDC details by name or ID.
func authenticateAndFetchSDDC(ctx context.Context, d *schema.ResourceData, client *Client) (SDDC, diag.Diagnostics) {
	// Extract input parameters.
	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	// Validate inputs.
	if sddcName == "" && sddcID == "" {
		return SDDC{}, diag.Errorf("neither 'sddc_name' nor 'sddc_id' was specified")
	}

	// Authenticate with VMware Cloud Services.
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		return SDDC{}, diag.Errorf("failed to authenticate with VMware Cloud Services: %s", err.Error())
	}

	// Authenticate the HCX API.
	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		return SDDC{}, diag.Errorf("failed to authenticate hcx api: %s", err.Error())
	}

	// Retrieve the SDDC details.
	var sddc SDDC
	if sddcID != "" {
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		sddc, err = GetSddcByName(client, sddcName)
	}
	if err != nil {
		return SDDC{}, diag.Errorf("failed to retrieve sddc details: %s", err.Error())
	}

	return sddc, nil
}

// calculateBackoff computes the exponential backoff duration with capped interval and jitter.
func calculateBackoff(retries int) time.Duration {
	base := constants.VmcRetryInterval * time.Duration(1<<retries)
	if base > constants.VmcMaxRetryInterval {
		base = constants.VmcMaxRetryInterval
	}
	// Add jitter to reduce the chance of simultaneous retries across instances.
	jitter := time.Duration(rand.Int63n(int64(base / 2)))
	return base + jitter
}

// waitWithContext waits for the backoff duration or respects context cancellations.
func waitWithContext(ctx context.Context, backoff time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(backoff):
		return nil
	}
}
