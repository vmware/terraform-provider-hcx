package main

import (
	"context"

	hcx "github.com/adeleporte/terraform-provider-hcx/hcx"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkBacking() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkBackingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vcuuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entityid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DistributedVirtualPortgroup",
			},
		},
	}
}

func dataSourceNetworkBackingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	network := d.Get("name").(string)
	vcuuid := d.Get("vcuuid").(string)
	network_type := d.Get("network_type").(string)

	res, err := hcx.GetNetworkBacking(client, vcuuid, network, network_type)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.EntityID)

	return diags
}
