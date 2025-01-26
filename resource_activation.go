// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-hcx/hcx"
)

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

func resourceActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	url := d.Get("url").(string)
	activationkey := d.Get("activationkey").(string)

	body := hcx.ActivateBody{
		Data: hcx.ActivateData{
			Items: []hcx.ActivateDataItem{
				{
					Config: hcx.ActivateDataItemConfig{
						URL:           url,
						ActivationKey: activationkey,
					},
				},
			},
		},
	}

	// First, check if already activated
	res, err := hcx.GetActivate(client)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Data.Items) == 0 {
		// No activation config found
		_, err := hcx.PostActivate(client, body)

		if err != nil {
			return diag.FromErr(err)
		}

		return resourceActivationRead(ctx, d, m)
	}

	d.SetId(res.Data.Items[0].Config.UUID)

	return resourceActivationRead(ctx, d, m)
}

func resourceActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	res, err := hcx.GetActivate(client)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Data.Items[0].Config.UUID)

	return diags
}

func resourceActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceActivationRead(ctx, d, m)
}

func resourceActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
