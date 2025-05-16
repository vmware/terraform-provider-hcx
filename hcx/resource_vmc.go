// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"time"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	tflog.Info(ctx, "Creating VMC resource", map[string]interface{}{
		"sddc_name": sddcName,
		"sddc_id":   sddcID,
	})

	// Authenticate with VMware Cloud Services
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		tflog.Error(ctx, "Error authenticating with VMware Cloud Services", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		tflog.Error(ctx, "Error authenticating with Cloud", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	var sddc SDDC
	if sddcID != "" {
		tflog.Debug(ctx, "Getting SDDC by ID", map[string]interface{}{
			"sddc_id": sddcID,
		})
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		tflog.Debug(ctx, "Getting SDDC by name", map[string]interface{}{
			"sddc_name": sddcName,
		})
		sddc, err = GetSddcByName(client, sddcName)
	}

	if err != nil {
		tflog.Error(ctx, "Error retrieving SDDC", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	// Check if already activated.
	if sddc.DeploymentStatus == constants.VmcActivationActiveStatus {
		tflog.Warn(ctx, "SDDC already activated", map[string]interface{}{
			"sddc_id": sddc.ID,
		})
		return diag.Errorf("Already activated")
	}

	// Activate HCX.
	tflog.Info(ctx, "Activating HCX on SDDC", map[string]interface{}{
		"sddc_id": sddc.ID,
	})
	_, err = ActivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		tflog.Error(ctx, "Error activating HCX on SDDC", map[string]interface{}{
			"error":   err.Error(),
			"sddc_id": sddc.ID,
		})
		return diag.FromErr(err)
	}

	// Wait for task to be completed.
	errcount := 0
	for {
		if sddcID != "" {
			sddc, err = GetSddcByID(client, sddcID)
		} else {
			sddc, err = GetSddcByName(client, sddcName)
		}
		if err != nil {
			// Attempt to bypass recurring situation where the HCX API
			// returns status 502 with a proxy server error, and an HTML response
			// instead of JSON.
			errcount++
			tflog.Warn(ctx, "Error retrieving SDDC status during activation", map[string]interface{}{
				"error":     err.Error(),
				"err_count": errcount,
			})
			if errcount > 12 {
				tflog.Error(ctx, "Max retries exceeded while checking SDDC activation status", map[string]interface{}{
					"error": err.Error(),
				})
				return diag.FromErr(err)
			}
		}

		if sddc.DeploymentStatus == constants.VmcActivationActiveStatus {
			tflog.Info(ctx, "SDDC activation completed successfully", map[string]interface{}{
				"sddc_id": sddc.ID,
			})
			break
		}

		if sddc.DeploymentStatus == constants.VmcActivationFailedStatus {
			tflog.Error(ctx, "SDDC activation failed", map[string]interface{}{
				"sddc_id": sddc.ID,
			})
			return diag.Errorf("Activation failed")
		}

		tflog.Debug(ctx, "Waiting for SDDC activation to complete", map[string]interface{}{
			"sddc_id":           sddc.ID,
			"deployment_status": sddc.DeploymentStatus,
		})
		time.Sleep(constants.VmcRetryInterval)
	}

	return resourceVmcRead(ctx, d, m)
}

// resourceVmcRead retrieves the VMware Cloud on AWS configuration.
func resourceVmcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	tflog.Debug(ctx, "Reading VMC resource", map[string]interface{}{
		"sddc_name": sddcName,
		"sddc_id":   sddcID,
	})

	if sddcName == "" && sddcID == "" {
		tflog.Error(ctx, "SDDC name or ID must be specified")
		return diag.Errorf("SDDC name or Id must be specified")
	}

	// Authenticate with VMware Cloud Services
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		tflog.Error(ctx, "Error authenticating with VMware Cloud Services", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		tflog.Error(ctx, "Error authenticating with Cloud", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Successfully authenticated with VMware Cloud Services", map[string]interface{}{
		"client": client.HostURL,
	})

	var sddc SDDC
	if sddcID != "" {
		tflog.Debug(ctx, "Getting SDDC by ID", map[string]interface{}{
			"sddc_id": sddcID,
		})
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		tflog.Debug(ctx, "Getting SDDC by name", map[string]interface{}{
			"sddc_name": sddcName,
		})
		sddc, err = GetSddcByName(client, sddcName)
	}
	if err != nil {
		tflog.Error(ctx, "Error retrieving SDDC", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	d.SetId(sddc.ID)
	if err := d.Set("cloud_url", sddc.CloudURL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_name", sddc.CloudName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_type", sddc.CloudType); err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Successfully read VMC resource", map[string]interface{}{
		"sddc_id":    sddc.ID,
		"cloud_url":  sddc.CloudURL,
		"cloud_name": sddc.CloudName,
		"cloud_type": sddc.CloudType,
	})

	return diags
}

// resourceVmcUpdate updates the VMware Cloud on AWS resource configuration.
func resourceVmcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Updating VMC resource", map[string]interface{}{
		"id": d.Id(),
	})
	return resourceVmcRead(ctx, d, m)
}

// resourceVmcDelete removes the VMware Cloud on AWS configuration and clears the state of the resource in the schema.
func resourceVmcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	tflog.Info(ctx, "Deleting VMC resource", map[string]interface{}{
		"sddc_name": sddcName,
		"sddc_id":   sddcID,
	})

	// Authenticate with VMware Cloud Services
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		tflog.Error(ctx, "Error authenticating with VMware Cloud Services", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		tflog.Error(ctx, "Error authenticating with Cloud", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	var sddc SDDC
	if sddcID != "" {
		tflog.Debug(ctx, "Getting SDDC by ID", map[string]interface{}{
			"sddc_id": sddcID,
		})
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		tflog.Debug(ctx, "Getting SDDC by name", map[string]interface{}{
			"sddc_name": sddcName,
		})
		sddc, err = GetSddcByName(client, sddcName)
	}
	if err != nil {
		tflog.Error(ctx, "Error retrieving SDDC", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	// Deactivate HCX
	tflog.Info(ctx, "Deactivating HCX on SDDC", map[string]interface{}{
		"sddc_id": sddc.ID,
	})
	_, err = DeactivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		tflog.Error(ctx, "Error deactivating HCX on SDDC", map[string]interface{}{
			"error":   err.Error(),
			"sddc_id": sddc.ID,
		})
		return diag.FromErr(err)
	}

	// Wait for task to be completed
	errcount := 0
	for {
		var sddc SDDC
		if sddcID != "" {
			sddc, err = GetSddcByID(client, sddcID)
		} else {
			sddc, err = GetSddcByName(client, sddcName)
		}
		if err != nil {
			// Attempt to bypass recurring situation where the HCX API
			// returns status 502 with a proxy server error, and an HTML response
			// instead of JSON.
			errcount++
			tflog.Warn(ctx, "Error retrieving SDDC status during deactivation", map[string]interface{}{
				"error":     err.Error(),
				"err_count": errcount,
			})
			if errcount > constants.VmcMaxRetries {
				tflog.Error(ctx, "Max retries exceeded while checking SDDC deactivation status", map[string]interface{}{
					"error": err.Error(),
				})
				return diag.FromErr(err)
			}
		}

		if sddc.DeploymentStatus == constants.VmcDeactivationInactiveStatus {
			tflog.Info(ctx, "SDDC deactivation completed successfully", map[string]interface{}{
				"sddc_id": sddc.ID,
			})
			break
		}

		if sddc.DeploymentStatus == constants.VmcDeactivationFailedStatus {
			tflog.Error(ctx, "SDDC deactivation failed", map[string]interface{}{
				"sddc_id": sddc.ID,
			})
			return diag.Errorf("Deactivation failed")
		}

		tflog.Debug(ctx, "Waiting for SDDC deactivation to complete", map[string]interface{}{
			"sddc_id":           sddc.ID,
			"deployment_status": sddc.DeploymentStatus,
		})
		time.Sleep(constants.VmcRetryInterval)
	}

	return diags
}
