package main

import (
	"context"

	hcx "github.com/adeleporte/terraform-provider-hcx/hcx"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceComputeProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceComputeProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	res, err := hcx.GetLocalCloudList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	network := d.Get("name").(string)

	cp, err := hcx.GetComputeProfile(client, res.Data.Items[0].EndpointId, network)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cp.ComputeProfileId)

	return diags
}
