// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hcx "github.com/vmware/terraform-provider-hcx/hcx"
)

func NetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vmc": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"mtu": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"prefix_length": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"gateway": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"site_pairing": {
			Type:     schema.TypeMap,
			Required: true,
		},
		"primary_dns": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"secondary_dns": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"dns_suffix": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"network_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "DistributedVirtualPortgroup",
		},
		"ip_range": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_address": {
						Type:     schema.TypeString,
						Required: true,
					},
					"end_address": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func resourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkProfileCreate,
		ReadContext:   resourceNetworkProfileRead,
		UpdateContext: resourceNetworkProfileUpdate,
		DeleteContext: resourceNetworkProfileDelete,

		Schema: NetSchema(),
	}
}

func resourceNetworkProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	vmc := d.Get("vmc").(bool)
	if vmc {
		// Dont create the network profile, just update it
		return resourceNetworkProfileUpdate(ctx, d, m)
	}

	mtu := d.Get("mtu").(int)
	prefix_length := d.Get("prefix_length").(int)

	name := d.Get("name").(string)
	gateway := d.Get("gateway").(string)

	primary_dns := d.Get("primary_dns").(string)
	secondary_dns := d.Get("secondary_dns").(string)
	dns_suffix := d.Get("dns_suffix").(string)

	sp := d.Get("site_pairing").(map[string]interface{})
	vcuuid := sp["local_vc"].(string)
	vclocalendpointid := sp["local_endpoint_id"].(string)

	network_name, ok := d.GetOk("network_name")
	if !ok && !vmc {
		return diag.Errorf("VMC switch is not enabled. Network name is mandatory")
	}
	network_type := d.Get("network_type").(string)
	network_id, err := hcx.GetNetworkBacking(client, vclocalendpointid, network_name.(string), network_type)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get IP Ranges from schema
	ip_range := d.Get("ip_range").([]interface{})

	ipr := []hcx.NetworkIpRange{}
	for _, j := range ip_range {
		s := j.(map[string]interface{})
		start_address := s["start_address"].(string)
		end_address := s["end_address"].(string)

		ipr = append(ipr, hcx.NetworkIpRange{
			StartAddress: start_address,
			EndAddress:   end_address,
		})
	}

	body := hcx.NetworkProfileBody{
		Name:         name,
		Organization: "DEFAULT",
		MTU:          mtu,
		Backings: []hcx.Backing{{
			BackingID:           network_id.EntityID,
			BackingName:         network_name.(string),
			VCenterInstanceUuid: vcuuid,
			Type:                network_type,
		},
		},
		IPScopes: []hcx.IPScope{
			{
				DnsSuffix:       dns_suffix,
				Gateway:         gateway,
				PrefixLength:    prefix_length,
				PrimaryDns:      primary_dns,
				SecondaryDns:    secondary_dns,
				NetworkIpRanges: ipr,
			},
		},
		L3TenantManaged: false,
		OwnedBySystem:   true,
	}

	res, err := hcx.InsertNetworkProfile(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for job completion
	for {
		jr, err := hcx.GetJobResult(client, res.Data.JobID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.IsDone {
			break
		}
		time.Sleep(5 * time.Second)
	}

	return resourceNetworkProfileRead(ctx, d, m)
}

func resourceNetworkProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)
	name := d.Get("name").(string)

	np, err := hcx.GetNetworkProfile(client, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(np.ObjectId)

	return diags
}

func resourceNetworkProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	// Get values from schema
	vmc := d.Get("vmc").(bool)
	mtu := d.Get("mtu").(int)
	prefix_length := d.Get("prefix_length").(int)

	name := d.Get("name").(string)
	gateway := d.Get("gateway").(string)

	primary_dns := d.Get("primary_dns").(string)
	secondary_dns := d.Get("secondary_dns").(string)
	dns_suffix := d.Get("dns_suffix").(string)
	network_name := d.Get("network_name").(string)
	network_type := d.Get("network_type").(string)

	sp := d.Get("site_pairing").(map[string]interface{})
	vcuuid := sp["local_vc"].(string)
	vclocalendpointid := sp["local_endpoint_id"].(string)

	// Get IP Ranges from schema
	ip_range := d.Get("ip_range").([]interface{})

	ipr := []hcx.NetworkIpRange{}
	for _, j := range ip_range {
		s := j.(map[string]interface{})
		start_address := s["start_address"].(string)
		end_address := s["end_address"].(string)

		ipr = append(ipr, hcx.NetworkIpRange{
			StartAddress: start_address,
			EndAddress:   end_address,
		})
	}

	// Read the exisint profile
	body, err := hcx.GetNetworkProfile(client, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the network profile

	if !vmc {
		body.Name = name

		// Get network details
		network_id, err := hcx.GetNetworkBacking(client, vclocalendpointid, network_name, network_type)
		if err != nil {
			return diag.FromErr(err)
		}

		body.Backings = []hcx.Backing{{
			BackingID:           network_id.EntityID,
			BackingName:         network_name,
			VCenterInstanceUuid: vcuuid,
			Type:                network_type,
		}}
	}

	body.MTU = mtu

	body.IPScopes = []hcx.IPScope{
		{
			DnsSuffix:       dns_suffix,
			Gateway:         gateway,
			PrefixLength:    prefix_length,
			PrimaryDns:      primary_dns,
			SecondaryDns:    secondary_dns,
			NetworkIpRanges: ipr,
			PoolID:          body.IPScopes[0].PoolID,
		},
	}

	res, err := hcx.UpdateNetworkProfile(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for job completion
	for {
		jr, err := hcx.GetJobResult(client, res.Data.JobID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.IsDone {
			break
		}
		time.Sleep(5 * time.Second)
	}

	return resourceNetworkProfileRead(ctx, d, m)
}

func resourceNetworkProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var res hcx.NetworkProfileResult
	var err error

	client := m.(*hcx.Client)
	//name := d.Get("name").(string)
	vmc := d.Get("vmc").(bool)

	if vmc {
		// If VMC, don't really delete the network profile
		// Read the exisint profile
		/*
			body, err := hcx.GetNetworkProfile(client, name)
			if err != nil {
				return diag.FromErr(err)
			}

			// Empty the IP Ranges
			body.IPScopes[0].NetworkIpRanges = []hcx.NetworkIpRange{}

			res, err = hcx.UpdateNetworkProfile(client, body)

			if err != nil {
				return diag.FromErr(err)
			}
		*/
		return diags
	} else {
		res, err = hcx.DeleteNetworkProfile(client, d.Id())

		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Wait for job completion
	for {
		jr, err := hcx.GetJobResult(client, res.Data.JobID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.IsDone {
			break
		}
		time.Sleep(5 * time.Second)
	}

	return diags
}
