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

func dataSourceNetworkBacking() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkBackingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the network backing.",
				Required:    true,
			},
			"vcuuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the vCenter instance associated with the network backing.",
				Required:    true,
			},
			"entityid": {
				Type:        schema.TypeString,
				Description: "The entity ID of the network backing.",
				Computed:    true,
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "The type of the network backing. Allowed values are 'DistributedVirtualPortgroup' and 'NsxtSegment'.",
				Optional:    true,
				Default:     "DistributedVirtualPortgroup",
			},
		},
	}
}

func dataSourceNetworkBackingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	network := d.Get("name").(string)
	vcUUID := d.Get("vcuuid").(string)
	networkType := d.Get("network_type").(string)

	res, err := hcx.GetNetworkBacking(client, vcUUID, network, networkType)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.EntityID)

	return diags
}
