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

func resourceSSOCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	url := d.Get("url").(string)

	body := hcx.InsertSSOBody{
		Data: hcx.InsertSSOData{
			Items: []hcx.InsertSSODataItem{
				{
					Config: hcx.InsertSSODataItemConfig{
						LookupServiceURL: url,
						ProviderType:     "PSC",
					},
				},
			},
		},
	}

	// First, check if SSO config is already present
	res, err := hcx.GetSSO(client)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.InsertSSOData.Items) == 0 {
		// No SSO config found
		res, err := hcx.InsertSSO(client, body)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(res.InsertSSOData.Items[0].Config.UUID)
		return resourceSSORead(ctx, d, m)
	}

	// Update existing SSO
	d.SetId(res.InsertSSOData.Items[0].Config.UUID)
	return resourceSSOUpdate(ctx, d, m)

}

func resourceSSORead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceSSOUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	url := d.Get("url").(string)

	body := hcx.InsertSSOBody{
		Data: hcx.InsertSSOData{
			Items: []hcx.InsertSSODataItem{
				{
					Config: hcx.InsertSSODataItemConfig{
						LookupServiceURL: url,
						UUID:             d.Id(),
						ProviderType:     "PSC",
					},
				},
			},
		},
	}

	_, err := hcx.UpdateSSO(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSSORead(ctx, d, m)
}

func resourceSSODelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	_, err := hcx.DeleteSSO(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
