// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceServiceMesh defines the resource schema for managing service mesh configuration.
func resourceServiceMesh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceMeshCreate,
		ReadContext:   resourceServiceMeshRead,
		UpdateContext: resourceServiceMeshUpdate,
		DeleteContext: resourceServiceMeshDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the service mesh.",
				Required:    true,
			},
			"local_compute_profile": {
				Type:        schema.TypeString,
				Description: "The local compute profile name.",
				Required:    true,
			},
			"remote_compute_profile": {
				Type:        schema.TypeString,
				Description: "The remote compute profile name.",
				Required:    true,
			},
			"app_path_resiliency_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable the Application Path Resiliency feature.",
				Optional:    true,
				Default:     false,
			},
			"tcp_flow_conditioning_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable the TCP flow conditioning feature.",
				Optional:    true,
				Default:     false,
			},
			"uplink_max_bandwidth": {
				Type:        schema.TypeInt,
				Description: "The maximum bandwidth used for uplinks.",
				Optional:    true,
				Default:     10000,
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Description: "Force delete of the service mesh. Sometimes needed when site pairing is no longer connected.",
				Optional:    true,
				Default:     false,
			},
			"service": {
				Type:        schema.TypeList,
				Description: "The list of HCX services.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the HCX service.",
							Required:    true,
						},
					},
				},
			},
			"site_pairing": {
				Type:        schema.TypeMap,
				Description: "The site pairing used by this service mesh.",
				Required:    true,
			},
			"nb_appliances": {
				Type:        schema.TypeInt,
				Description: "The number of Network Extension appliances to deploy.",
				Optional:    true,
				Default:     1,
			},
			"appliances_id": {
				Type:        schema.TypeList,
				Description: "The IDs of the Network Extension appliances.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the Network Extension appliance.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// resourceServiceMeshCreate creates the service mesh configuration.
func resourceServiceMeshCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	name := d.Get("name").(string)
	sitePairing := d.Get("site_pairing").(map[string]interface{})
	localEndpointID := sitePairing["local_endpoint_id"].(string)
	localEndpointName := sitePairing["local_name"].(string)

	remoteEndpointID := sitePairing["id"].(string)
	remoteEndpointName := sitePairing["remote_name"].(string)

	uplinkMaxBandwidth := d.Get("uplink_max_bandwidth").(int)
	appPathResiliencyEnabled := d.Get("app_path_resiliency_enabled").(bool)
	tcpFlowConditioningEnabled := d.Get("tcp_flow_conditioning_enabled").(bool)

	services := d.Get("service").([]interface{})
	servicesFromSchema := []Service{}
	for _, j := range services {
		s := j.(map[string]interface{})
		name := s["name"].(string)

		sTmp := Service{
			Name: name,
		}
		servicesFromSchema = append(servicesFromSchema, sTmp)
	}

	remoteComputeProfileName := d.Get("remote_compute_profile").(string)
	remoteComputeProfile, err := GetComputeProfile(client, remoteEndpointID, remoteComputeProfileName)
	if err != nil {
		return diag.FromErr(err)
	}

	localComputeProfileName := d.Get("local_compute_profile").(string)
	localComputeProfile, err := GetComputeProfile(client, localEndpointID, localComputeProfileName)
	if err != nil {
		return diag.FromErr(err)
	}

	nbAppliances := d.Get("nb_appliances").(int)

	body := InsertServiceMeshBody{
		Name: name,
		ComputeProfiles: []ComputeProfile{
			{
				EndpointID:         localEndpointID,
				EndpointName:       localEndpointName,
				ComputeProfileID:   localComputeProfile.ComputeProfileID,
				ComputeProfileName: localComputeProfile.Name,
			},
			{
				EndpointID:         remoteEndpointID,
				EndpointName:       remoteEndpointName,
				ComputeProfileID:   remoteComputeProfile.ComputeProfileID,
				ComputeProfileName: remoteComputeProfile.Name,
			},
		},
		WanoptConfig: WanoptConfig{
			UplinkMaxBandwidth: uplinkMaxBandwidth,
		},
		TrafficEnggCfg: TrafficEnggCfg{
			IsAppPathResiliencyEnabled:   appPathResiliencyEnabled,
			IsTCPFlowConditioningEnabled: tcpFlowConditioningEnabled,
		},
		Services: servicesFromSchema,
		SwitchPairCount: []SwitchPairCount{
			{
				Switches: []Switch{
					localComputeProfile.Switches[0],
					remoteComputeProfile.Switches[0],
				},
				L2cApplianceCount: nbAppliances,
			},
		},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return diag.FromErr(err)
	}

	res2, err := InsertServiceMesh(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task completion
	for {
		jr, err := GetTaskResult(client, res2.Data.InterconnectID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Status == constants.SuccessStatus {
			break
		}

		if jr.Status == constants.FailedStatus {
			return diag.FromErr(errors.New("task failed"))
		}

		time.Sleep(5 * time.Second)
	}

	// Update Appliances ID
	appliances, err := GetAppliances(client, sitePairing["local_endpoint_id"].(string), res2.Data.ServiceMeshID)
	if err != nil {
		return diag.FromErr(err)
	}

	tmp := []map[string]string{}

	for _, j := range appliances {
		a := map[string]string{}
		a["id"] = j.ApplianceID
		tmp = append(tmp, a)
	}
	d.Set("appliances_id", tmp)

	d.SetId(res2.Data.ServiceMeshID)

	return resourceServiceMeshRead(ctx, d, m)

}

// resourceServiceMeshRead retrieves the service mesh configuration.
func resourceServiceMeshRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

// resourceServiceMeshUpdate updates the service mesh configuration.
func resourceServiceMeshUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceServiceMeshRead(ctx, d, m)
}

// resourceServiceMeshDelete removes the service mesh configuration and clears the state of the resource in the schema.
func resourceServiceMeshDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)
	force := d.Get("force_delete").(bool)

	res, err := DeleteServiceMesh(client, d.Id(), force)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task completion
	for {
		jr, err := GetTaskResult(client, res.Data.InterconnectTaskID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Status == constants.SuccessStatus {
			break
		}

		if jr.Status == constants.FailedStatus {
			return diag.FromErr(errors.New("task failed"))
		}

		time.Sleep(5 * time.Second)
	}

	return diags
}
