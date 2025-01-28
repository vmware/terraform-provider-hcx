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

// dataSourceComputeProfile defines the data source schema to retrieve information about a compute profile.
func dataSourceComputeProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the compute profile.",
				Required:    true,
			},
			"vcenter": {
				Type:        schema.TypeString,
				Description: "The ID of the local vCenter instance.",
				Required:    true,
			},
		},
	}
}

// dataSourceComputeProfileRead retrieves the compute profile configuration.
func dataSourceComputeProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	res, err := hcx.GetLocalCloudList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	network := d.Get("name").(string)

	cp, err := hcx.GetComputeProfile(client, res.Data.Items[0].EndpointID, network)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cp.ComputeProfileID)

	return diags
}
