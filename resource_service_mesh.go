// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hcx "github.com/vmware/terraform-provider-hcx/hcx"
)

func resourceServiceMesh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceMeshCreate,
		ReadContext:   resourceServiceMeshRead,
		UpdateContext: resourceServiceMeshUpdate,
		DeleteContext: resourceServiceMeshDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_compute_profile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_compute_profile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_path_resiliency_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tcp_flow_conditioning_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"uplink_max_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10000,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service": {
				Type:        schema.TypeList,
				Description: "Rules description",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"site_pairing": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"nb_appliances": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"appliances_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceServiceMeshCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	name := d.Get("name").(string)
	site_pairing := d.Get("site_pairing").(map[string]interface{})
	local_endpoint_id := site_pairing["local_endpoint_id"].(string)
	local_endpoint_name := site_pairing["local_name"].(string)

	remote_endpoint_id := site_pairing["id"].(string)
	remote_endpoint_name := site_pairing["remote_name"].(string)

	uplink_max_bandwidth := d.Get("uplink_max_bandwidth").(int)
	app_path_resiliency_enabled := d.Get("app_path_resiliency_enabled").(bool)
	tcp_flow_conditioning_enabled := d.Get("tcp_flow_conditioning_enabled").(bool)

	services := d.Get("service").([]interface{})
	services_from_schema := []hcx.Service{}
	for _, j := range services {
		s := j.(map[string]interface{})
		name := s["name"].(string)

		s_tmp := hcx.Service{
			Name: name,
		}
		services_from_schema = append(services_from_schema, s_tmp)
	}

	remote_compute_profile_name := d.Get("remote_compute_profile").(string)
	remote_compute_profile, err := hcx.GetComputeProfile(client, remote_endpoint_id, remote_compute_profile_name)
	if err != nil {
		return diag.FromErr(err)
	}

	local_compute_profile_name := d.Get("local_compute_profile").(string)
	local_compute_profile, err := hcx.GetComputeProfile(client, local_endpoint_id, local_compute_profile_name)
	if err != nil {
		return diag.FromErr(err)
	}

	nb_appliances := d.Get("nb_appliances").(int)

	body := hcx.InsertServiceMeshBody{
		Name: name,
		ComputeProfiles: []hcx.ComputeProfile{
			{
				EndpointId:         local_endpoint_id,
				EndpointName:       local_endpoint_name,
				ComputeProfileId:   local_compute_profile.ComputeProfileId,
				ComputeProfileName: local_compute_profile.Name,
			},
			{
				EndpointId:         remote_endpoint_id,
				EndpointName:       remote_endpoint_name,
				ComputeProfileId:   remote_compute_profile.ComputeProfileId,
				ComputeProfileName: remote_compute_profile.Name,
			},
		},
		WanoptConfig: hcx.WanoptConfig{
			UplinkMaxBandwidth: uplink_max_bandwidth,
		},
		TrafficEnggCfg: hcx.TrafficEnggCfg{
			IsAppPathResiliencyEnabled:   app_path_resiliency_enabled,
			IsTcpFlowConditioningEnabled: tcp_flow_conditioning_enabled,
		},
		Services: services_from_schema,
		SwitchPairCount: []hcx.SwitchPairCount{
			{
				Switches: []hcx.Switch{
					local_compute_profile.Switches[0],
					remote_compute_profile.Switches[0],
				},
				L2cApplianceCount: nb_appliances,
			},
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	res2, err := hcx.InsertServiceMesh(client, body)

	if err != nil {
		//return diag.FromErr(errors.New(fmt.Sprintf("%s", buf)))
		return diag.FromErr(err)
	}

	// Wait for task completion
	for {
		jr, err := hcx.GetTaskResult(client, res2.Data.InterconnectTaskId)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Status == "SUCCESS" {
			break
		}

		if jr.Status == "FAILED" {
			return diag.FromErr(errors.New("task failed"))
		}

		time.Sleep(5 * time.Second)
	}

	// Update Appliances ID
	appliances, err := hcx.GetAppliances(client, site_pairing["local_endpoint_id"].(string), res2.Data.ServiceMeshId)
	if err != nil {
		return diag.FromErr(err)
	}

	tmp := []map[string]string{}

	for _, j := range appliances {
		a := map[string]string{}
		a["id"] = j.ApplianceId
		tmp = append(tmp, a)
	}
	d.Set("appliances_id", tmp)

	d.SetId(res2.Data.ServiceMeshId)

	return resourceServiceMeshRead(ctx, d, m)

}

func resourceServiceMeshRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceServiceMeshUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceServiceMeshRead(ctx, d, m)
}

func resourceServiceMeshDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)
	force := d.Get("force_delete").(bool)

	res, err := hcx.DeleteServiceMesh(client, d.Id(), force)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task completion
	for {
		jr, err := hcx.GetTaskResult(client, res.Data.InterconnectTaskId)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Status == "SUCCESS" {
			break
		}

		if jr.Status == "FAILED" {
			return diag.FromErr(errors.New("task failed"))
		}

		time.Sleep(5 * time.Second)
	}

	return diags
}
