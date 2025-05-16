// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"fmt"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"
	"github.com/vmware/terraform-provider-hcx/hcx/validators"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceNetworkBacking defines a data source schema to retrieve information about a network backing.
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
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("The network type for the network backing. Allowed values include: %v.", constants.AllowedNetworkTypes),
				Optional:     true,
				Default:      constants.NetworkTypeDvpg,
				ValidateFunc: validators.ValidateNetworkType,
			},
		},
	}
}

// dataSourceNetworkBackingRead retrieves the network backing configuration.
func dataSourceNetworkBackingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Reading network backing data source")

	client := m.(*Client)

	network := d.Get("name").(string)
	vcUUID := d.Get("vcuuid").(string)
	networkType := d.Get("network_type").(string)

	tflog.Debug(ctx, "Getting network backing", map[string]interface{}{
		"name":         network,
		"vcuuid":       vcUUID,
		"network_type": networkType,
	})

	res, err := GetNetworkBacking(client, vcUUID, network, networkType)
	if err != nil {
		tflog.Error(ctx, "Failed to get network backing", map[string]interface{}{
			"error":        err.Error(),
			"network":      network,
			"network_type": networkType,
		})
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Found network backing", map[string]interface{}{
		"entity_id": res.EntityID,
		"name":      network,
	})
	d.SetId(res.EntityID)

	return diags
}
