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

// resourceSSO defines the resource schema for managing SSO configuration.
func resourceSSO() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSOCreate,
		ReadContext:   resourceSSORead,
		UpdateContext: resourceSSOUpdate,
		DeleteContext: resourceSSODelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The URL of the vCenter instance.",
				Required:    true,
			},
			"vcenter": {
				Type:        schema.TypeString,
				Description: "The ID of the vCenter instance.",
				Required:    true,
			},
		},
	}
}

// resourceSSOCreate creates the SSO configuration.
func resourceSSOCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	url := d.Get("url").(string)

	tflog.Info(ctx, "Creating SSO configuration", map[string]interface{}{
		"url": url,
	})

	body := InsertSSOBody{
		Data: InsertSSOData{
			Items: []InsertSSODataItem{
				{
					Config: InsertSSODataItemConfig{
						LookupServiceURL: url,
						ProviderType:     constants.DefaultSsoProviderType,
					},
				},
			},
		},
	}

	// First, check if SSO config is already present
	tflog.Debug(ctx, "Checking if SSO configuration already exists")
	res, err := GetSSO(client)
	if err != nil {
		tflog.Error(ctx, "Failed to get existing SSO configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	if len(res.InsertSSOData.Items) == 0 {
		// No SSO configuration found.
		tflog.Debug(ctx, "No existing SSO configuration found, creating new")
		res, err := InsertSSO(client, body)

		if err != nil {
			tflog.Error(ctx, "Failed to insert SSO configuration", map[string]interface{}{
				"error": err.Error(),
				"url":   url,
			})
			return diag.FromErr(err)
		}

		tflog.Info(ctx, "SSO configuration created successfully", map[string]interface{}{
			"uuid": res.InsertSSOData.Items[0].Config.UUID,
		})
		d.SetId(res.InsertSSOData.Items[0].Config.UUID)
		return resourceSSORead(ctx, d, m)
	}

	// Update existing SSO configuration.
	tflog.Debug(ctx, "Existing SSO configuration found, updating instead", map[string]interface{}{
		"uuid": res.InsertSSOData.Items[0].Config.UUID,
	})
	d.SetId(res.InsertSSOData.Items[0].Config.UUID)
	return resourceSSOUpdate(ctx, d, m)
}

// resourceSSORead retrieves the SSO configuration.
func resourceSSORead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tflog.Debug(ctx, "Reading SSO configuration", map[string]interface{}{
		"uuid": d.Id(),
	})
	return diags
}

// resourceSSOUpdate updates the SSO configuration.
func resourceSSOUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	url := d.Get("url").(string)

	tflog.Info(ctx, "Updating SSO configuration", map[string]interface{}{
		"uuid": d.Id(),
		"url":  url,
	})

	body := InsertSSOBody{
		Data: InsertSSOData{
			Items: []InsertSSODataItem{
				{
					Config: InsertSSODataItemConfig{
						LookupServiceURL: url,
						UUID:             d.Id(),
						ProviderType:     constants.DefaultSsoProviderType,
					},
				},
			},
		},
	}

	_, err := UpdateSSO(client, body)

	if err != nil {
		tflog.Error(ctx, "Failed to update SSO configuration", map[string]interface{}{
			"error": err.Error(),
			"uuid":  d.Id(),
			"url":   url,
		})
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "SSO configuration updated successfully", map[string]interface{}{
		"uuid": d.Id(),
	})
	return resourceSSORead(ctx, d, m)
}

// resourceSSODelete removes the SSO configuration and clears the state of the resource in the schema.
func resourceSSODelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	tflog.Info(ctx, "Deleting SSO configuration", map[string]interface{}{
		"uuid": d.Id(),
	})

	_, err := DeleteSSO(client, d.Id())
	if err != nil {
		tflog.Error(ctx, "Failed to delete SSO configuration", map[string]interface{}{
			"error": err.Error(),
			"uuid":  d.Id(),
		})
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "SSO configuration deleted successfully")
	return diags
}
