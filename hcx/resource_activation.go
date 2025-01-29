// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"

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
				Default:     "https://connect.hcx.vmware.com",
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

	client := m.(*Client)

	url := d.Get("url").(string)
	activationkey := d.Get("activationkey").(string)

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
	res, err := GetActivate(client)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Data.Items) == 0 {
		// No activation config found
		_, err := PostActivate(client, body)

		if err != nil {
			return diag.FromErr(err)
		}

		return resourceActivationRead(ctx, d, m)
	}

	d.SetId(res.Data.Items[0].Config.UUID)

	return resourceActivationRead(ctx, d, m)
}

// resourceActivationRead retrieves the activation configuration and sets the resource ID in the schema.
func resourceActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	res, err := GetActivate(client)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Data.Items[0].Config.UUID)

	return diags
}

// resourceActivationUpdate updates the activation configuration by invoking the read operation to refresh its state.
func resourceActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceActivationRead(ctx, d, m)
}

// resourceActivationDelete removes the activation configuration and clears the state of the resource in the schema.
func resourceActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
