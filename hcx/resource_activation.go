// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceActivation defines the resource schema for managing activation configurations.
func resourceActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationCreate,
		ReadContext:   resourceActivationRead,
		UpdateContext: resourceActivationUpdate,
		DeleteContext: resourceActivationDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The URL for activation.",
				Optional:    true,
				Default:     constants.HcxCloudURL,
			},
			"activationkey": {
				Type:        schema.TypeString,
				Description: "The activation key.",
				Required:    true,
			},
		},
	}
}

// resourceActivationCreate creates the activation configuration resource.
func resourceActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Creating activation")

	client := m.(*Client)

	url := d.Get("url").(string)
	activationkey := d.Get("activationkey").(string)

	tflog.Debug(ctx, "Activation parameters", map[string]interface{}{
		"url": url,
	})

	body := ActivateBody{
		Data: ActivateData{
			Items: []ActivateDataItem{
				{
					Config: ActivateDataItemConfig{
						URL:           url,
						ActivationKey: activationkey,
					},
				},
			},
		},
	}

	// First, check if already activated
	tflog.Debug(ctx, "Checking if already activated")
	res, err := GetActivate(client)
	if err != nil {
		tflog.Error(ctx, "Failed to get activation status", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	if len(res.Data.Items) == 0 {
		// No activation config found
		tflog.Info(ctx, "No activation found, activating HCX")
		_, err := PostActivate(client, body)

		if err != nil {
			tflog.Error(ctx, "Failed to activate HCX", map[string]interface{}{
				"error": err.Error(),
			})
			return diag.FromErr(err)
		}

		tflog.Info(ctx, "HCX activated successfully")
		return resourceActivationRead(ctx, d, m)
	}

	tflog.Info(ctx, "HCX already activated", map[string]interface{}{
		"uuid": res.Data.Items[0].Config.UUID,
	})
	d.SetId(res.Data.Items[0].Config.UUID)

	return resourceActivationRead(ctx, d, m)
}

// resourceActivationRead retrieves the activation configuration and sets the resource ID in the schema.
func resourceActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Reading activation configuration")

	client := m.(*Client)

	tflog.Debug(ctx, "Getting activation status")
	res, err := GetActivate(client)
	if err != nil {
		tflog.Error(ctx, "Failed to get activation status", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	if len(res.Data.Items) > 0 {
		tflog.Debug(ctx, "Found activation configuration", map[string]interface{}{
			"uuid": res.Data.Items[0].Config.UUID,
		})
		d.SetId(res.Data.Items[0].Config.UUID)
	} else {
		tflog.Warn(ctx, "No activation configuration found")
		d.SetId("")
	}

	return diags
}

// resourceActivationUpdate updates the activation configuration by invoking the read operation to refresh its state.
func resourceActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Updating activation", map[string]interface{}{
		"id": d.Id(),
	})

	return resourceActivationRead(ctx, d, m)
}

// resourceActivationDelete removes the activation configuration and clears the state of the resource in the schema.
func resourceActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Deleting activation", map[string]interface{}{
		"id": d.Id(),
	})

	// Note: HCX activation cannot actually be deleted via API
	tflog.Warn(ctx, "HCX activation cannot be deleted via API, only removing from Terraform state")

	return diags
}
