// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-hcx/hcx"
)

func resourceL2Extension() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL2ExtensionCreate,
		ReadContext:   resourceL2ExtensionRead,
		UpdateContext: resourceL2ExtensionUpdate,
		DeleteContext: resourceL2ExtensionDelete,

		Schema: map[string]*schema.Schema{
			"site_pairing": {
				Type:        schema.TypeMap,
				Description: "The site pairing used by this service mesh.",
				Required:    true,
			},
			"service_mesh_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Service Mesh to be used for this L2 extension.",
				Required:    true,
			},
			"source_network": {
				Type:        schema.TypeString,
				Description: "The source network. Must be a distributed port group which is VLAN tagged.",
				Required:    true,
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "The network backing type. Allowed values include: 'DistributedVirtualPortgroup' and 'NsxtSegment'.",
				Optional:    true,
				Default:     "DistributedVirtualPortgroup",
			},
			"destination_t1": {
				Type:        schema.TypeString,
				Description: "The name of the NSX T1 at the destination.",
				Required:    true,
			},
			"gateway": {
				Type:        schema.TypeString,
				Description: "The gateway address to configure on the NSX T1. Should be equal to the existing default gateway at the source site.",
				Optional:    true,
				Default:     false,
			},
			"netmask": {
				Type:        schema.TypeString,
				Description: "The netmask.",
				Optional:    true,
				Default:     false,
			},
			"mon": {
				Type:        schema.TypeBool,
				Description: "Enable the MON (Mobility Optimized Networking) feature.",
				Optional:    true,
				Default:     false,
			},
			"egress_optimization": {
				Type:        schema.TypeBool,
				Description: "Enable the Egress Optimization feature.",
				Optional:    true,
				Default:     false,
			},
			"appliance_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network Extension appliance to use for the L2 extension.",
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func resourceL2ExtensionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	sitePairing := d.Get("site_pairing").(map[string]interface{})
	vcGUID := sitePairing["local_vc"].(string)

	//service_mesh := d.Get("service_mesh").(map[string]interface{})
	sourceNetwork := d.Get("source_network").(string)
	destinationT1 := d.Get("destination_t1").(string)
	gateway := d.Get("gateway").(string)
	netmask := d.Get("netmask").(string)

	destinationEndpointID := sitePairing["id"].(string)
	destinationEndpointName := sitePairing["remote_name"].(string)
	destinationEndpointType := sitePairing["remote_endpoint_type"].(string)

	destinationResourceID := sitePairing["remote_resource_id"].(string)
	destinationResourceName := sitePairing["remote_resource_name"].(string)
	destinationResourceType := sitePairing["remote_resource_type"].(string)

	mon := d.Get("mon").(bool)
	egressOptimization := d.Get("egress_optimization").(bool)

	networkType := d.Get("network_type").(string)

	serviceMeshID := d.Get("service_mesh_id").(string)

	dvpg, err := hcx.GetNetworkBacking(client, sitePairing["local_endpoint_id"].(string), sourceNetwork, networkType)
	if err != nil {
		return diag.FromErr(err)
	}

	applianceID := d.Get("appliance_id").(string)
	if applianceID == "" {
		// GET THE FIRST APPLIANCE
		appliance, err := hcx.GetAppliance(client, sitePairing["local_endpoint_id"].(string), serviceMeshID)
		if err != nil {
			return diag.FromErr(err)
		}
		applianceID = appliance.ApplianceID
	}

	body := hcx.InsertL2ExtensionBody{
		Gateway: gateway,
		Netmask: netmask,
		DestinationNetwork: hcx.DestinationNetwork{
			GatewayID: destinationT1,
		},
		DNS: []string{},
		Features: hcx.Features{
			EgressOptimization: egressOptimization,
			Mon:                mon,
		},
		SourceAppliance: hcx.SourceAppliance{
			ApplianceID: applianceID,
		},
		SourceNetwork: hcx.SourceNetwork{
			NetworkID:   dvpg.EntityID,
			NetworkName: dvpg.Name,
			NetworkType: dvpg.EntityType,
		},
		VcGUID: vcGUID,
		Destination: hcx.Destination{
			EndpointID:   destinationEndpointID,
			EndpointName: destinationEndpointName,
			EndpointType: destinationEndpointType,
			ResourceID:   destinationResourceID,
			ResourceName: destinationResourceName,
			ResourceType: destinationResourceType,
		},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return diag.FromErr(err)
	}

	res2, err := hcx.InsertL2Extension(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for job completion
	for {
		jr, err := hcx.GetJobResult(client, res2.ID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.IsDone {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Get L2 Extension ID
	res3, err := hcx.GetL2Extensions(client, dvpg.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res3.StretchID)

	return resourceL2ExtensionRead(ctx, d, m)

}

func resourceL2ExtensionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceL2ExtensionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceL2ExtensionRead(ctx, d, m)
}

func resourceL2ExtensionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	res, err := hcx.DeleteL2Extension(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for job completion
	for {
		jr, err := hcx.GetJobResult(client, res.ID)
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
