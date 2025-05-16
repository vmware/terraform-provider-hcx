// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	tflog.Info(ctx, "Reading compute profile data source")

	client := m.(*Client)

	tflog.Debug(ctx, "Getting local cloud list")
	res, err := GetLocalCloudList(client)
	if err != nil {
		tflog.Error(ctx, "Failed to get local cloud list", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	networkName := d.Get("name").(string)
	tflog.Debug(ctx, "Looking for compute profile", map[string]interface{}{
		"name": networkName,
	})

	cp, err := GetComputeProfile(client, res.Data.Items[0].EndpointID, networkName)
	if err != nil {
		tflog.Error(ctx, "Failed to get compute profile", map[string]interface{}{
			"error": err.Error(),
			"name":  networkName,
		})
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Found compute profile", map[string]interface{}{
		"id":   cp.ComputeProfileID,
		"name": networkName,
	})
	d.SetId(cp.ComputeProfileID)

	return diags
}
