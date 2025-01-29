// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NetSchema defines the resource schema a network profile.
func NetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vmc": {
			Type:        schema.TypeBool,
			Description: "If set to true, the network profile will not be created or deleted, only IP pools will be updated.",
			Optional:    true,
			Default:     false,
		},
		"mtu": {
			Type:        schema.TypeInt,
			Description: "The MTU of the network profile.",
			Required:    true,
		},
		"prefix_length": {
			Type:        schema.TypeInt,
			Description: "The prefix length for the network profile.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the network profile.",
			Required:    true,
		},
		"gateway": {
			Type:        schema.TypeString,
			Description: "The gateway for the network profile.",
			Optional:    true,
			Default:     "",
		},
		"site_pairing": {
			Type:        schema.TypeMap,
			Description: "The site pairing map, to be retrieved with the 'hcx_site_pairing' resource.",
			Required:    true,
		},
		"primary_dns": {
			Type:        schema.TypeString,
			Description: "The primary DNS server for the network profile.",
			Optional:    true,
			Default:     "",
		},
		"secondary_dns": {
			Type:        schema.TypeString,
			Description: "The secondary DNS server for the network profile.",
			Optional:    true,
			Default:     "",
		},
		"dns_suffix": {
			Type:        schema.TypeString,
			Description: "The DNS suffix for the network profile.",
			Optional:    true,
			Default:     "",
		},
		"network_name": {
			Type:        schema.TypeString,
			Description: "The network name used for the network profile.",
			Optional:    true,
		},
		"network_type": {
			Type:        schema.TypeString,
			Description: "The network type for the network profile.",
			Optional:    true,
			Default:     "DistributedVirtualPortgroup",
		},
		"ip_range": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_address": {
						Type:        schema.TypeString,
						Description: "The start address of the IP pool for the network profile.",
						Required:    true,
					},
					"end_address": {
						Type:        schema.TypeString,
						Description: "The end address of the IP pool for the network profile.",
						Required:    true,
					},
				},
			},
		},
	}
}

// resourceComputeProfile defines the resource for managing network profile configuration.
func resourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkProfileCreate,
		ReadContext:   resourceNetworkProfileRead,
		UpdateContext: resourceNetworkProfileUpdate,
		DeleteContext: resourceNetworkProfileDelete,

		Schema: NetSchema(),
	}
}

// resourceNetworkProfileCreate creates the network profile configuration.
func resourceNetworkProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	vmc := d.Get("vmc").(bool)
	if vmc {
		// Dont create the network profile, just update it
		return resourceNetworkProfileUpdate(ctx, d, m)
	}

	mtu := d.Get("mtu").(int)
	prefixLength := d.Get("prefix_length").(int)

	name := d.Get("name").(string)
	gateway := d.Get("gateway").(string)

	primaryDNS := d.Get("primary_dns").(string)
	secondaryDNS := d.Get("secondary_dns").(string)
	dnsSuffix := d.Get("dns_suffix").(string)

	sp := d.Get("site_pairing").(map[string]interface{})
	vcUUID := sp["local_vc"].(string)
	vcLocalEndpointID := sp["local_endpoint_id"].(string)

	networkName, ok := d.GetOk("network_name")
	if !ok && !vmc {
		return diag.Errorf("VMC switch is not enabled. Network name is mandatory")
	}
	networkType := d.Get("network_type").(string)
	networkiD, err := GetNetworkBacking(client, vcLocalEndpointID, networkName.(string), networkType)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get IP Ranges from schema
	ipRange := d.Get("ip_range").([]interface{})

	ipr := []NetworkIPRange{}
	for _, j := range ipRange {
		s := j.(map[string]interface{})
		startAddress := s["start_address"].(string)
		endAddress := s["end_address"].(string)

		ipr = append(ipr, NetworkIPRange{
			StartAddress: startAddress,
			EndAddress:   endAddress,
		})
	}

	body := NetworkProfileBody{
		Name:         name,
		Organization: "DEFAULT",
		MTU:          mtu,
		Backings: []Backing{{
			BackingID:           networkiD.EntityID,
			BackingName:         networkName.(string),
			VCenterInstanceUUID: vcUUID,
			Type:                networkType,
		},
		},
		IPScopes: []IPScope{
			{
				DNSSuffix:       dnsSuffix,
				Gateway:         gateway,
				PrefixLength:    prefixLength,
				PrimaryDNS:      primaryDNS,
				SecondaryDNS:    secondaryDNS,
				NetworkIPRanges: ipr,
			},
		},
		L3TenantManaged: false,
		OwnedBySystem:   true,
	}

	res, err := InsertNetworkProfile(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for job completion
	for {
		jr, err := GetJobResult(client, res.Data.JobID)
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

// resourceNetworkProfileRead retrieves the network profile configuration.
func resourceNetworkProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)
	name := d.Get("name").(string)

	np, err := GetNetworkProfile(client, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(np.ObjectID)

	return diags
}

// resourceNetworkProfileUpdate updates the network profile configuration.
func resourceNetworkProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	// Get values from schema
	vmc := d.Get("vmc").(bool)
	mtu := d.Get("mtu").(int)
	prefixLength := d.Get("prefix_length").(int)

	name := d.Get("name").(string)
	gateway := d.Get("gateway").(string)

	primaryDNS := d.Get("primary_dns").(string)
	secondaryDNS := d.Get("secondary_dns").(string)
	dnsSuffix := d.Get("dns_suffix").(string)
	networkName := d.Get("network_name").(string)
	networkType := d.Get("network_type").(string)

	sp := d.Get("site_pairing").(map[string]interface{})
	vcUUID := sp["local_vc"].(string)
	vcLocalEndpointID := sp["local_endpoint_id"].(string)

	// Get IP Ranges from schema
	ipRange := d.Get("ip_range").([]interface{})

	ipr := []NetworkIPRange{}
	for _, j := range ipRange {
		s := j.(map[string]interface{})
		startAddress := s["start_address"].(string)
		endAddress := s["end_address"].(string)

		ipr = append(ipr, NetworkIPRange{
			StartAddress: startAddress,
			EndAddress:   endAddress,
		})
	}

	// Read the existing profile
	body, err := GetNetworkProfile(client, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the network profile

	if !vmc {
		body.Name = name

		// Get network details
		networkID, err := GetNetworkBacking(client, vcLocalEndpointID, networkName, networkType)
		if err != nil {
			return diag.FromErr(err)
		}

		body.Backings = []Backing{{
			BackingID:           networkID.EntityID,
			BackingName:         networkName,
			VCenterInstanceUUID: vcUUID,
			Type:                networkType,
		}}
	}

	body.MTU = mtu

	body.IPScopes = []IPScope{
		{
			DNSSuffix:       dnsSuffix,
			Gateway:         gateway,
			PrefixLength:    prefixLength,
			PrimaryDNS:      primaryDNS,
			SecondaryDNS:    secondaryDNS,
			NetworkIPRanges: ipr,
			PoolID:          body.IPScopes[0].PoolID,
		},
	}

	res, err := UpdateNetworkProfile(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for job completion
	for {
		jr, err := GetJobResult(client, res.Data.JobID)
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

// resourceNetworkProfileDelete removes the network profile configuration and clears the state of the resource in the
// schema.
func resourceNetworkProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var res NetworkProfileResult
	var err error

	client := m.(*Client)
	//name := d.Get("name").(string)
	vmc := d.Get("vmc").(bool)

	if vmc {
		// If VMware Cloud on AWS, don't really delete the network profile
		// Read the existing profile
		/*
			body, err := hcx.GetNetworkProfile(client, name)
			if err != nil {
				return diag.FromErr(err)
			}

			// Empty the IP Ranges
			body.IPScopes[0].NetworkIPRanges = []hcx.NetworkIPRange{}

			res, err = hcx.UpdateNetworkProfile(client, body)

			if err != nil {
				return diag.FromErr(err)
			}
		*/
		return diags
	} else {
		res, err = DeleteNetworkProfile(client, d.Id())

		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Wait for job completion
	for {
		jr, err := GetJobResult(client, res.Data.JobID)
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
