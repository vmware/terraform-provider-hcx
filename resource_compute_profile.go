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
	"github.com/vmware/terraform-provider-hcx/hcx"
)

func resourceComputeProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeProfileCreate,
		ReadContext:   resourceComputeProfileRead,
		UpdateContext: resourceComputeProfileUpdate,
		DeleteContext: resourceComputeProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the compute profile.",
				Required:    true,
			},
			"datacenter": {
				Type:        schema.TypeString,
				Description: "The datacenter where HCX services will be available.",
				Required:    true,
			},
			"cluster": {
				Type:        schema.TypeString,
				Description: "The cluster used for HCX appliances deployment.",
				Required:    true,
			},
			"datastore": {
				Type:        schema.TypeString,
				Description: "The datastore used for HCX appliances deployment.",
				Optional:    true,
				Default:     "",
			},
			"management_network": {
				Type:        schema.TypeString,
				Description: "The management network profile (ID).",
				Required:    true,
			},
			"replication_network": {
				Type:        schema.TypeString,
				Description: "The replication network profile (ID).",
				Optional:    true,
				Default:     "",
			},
			"uplink_network": {
				Type:        schema.TypeString,
				Description: "The uplink network profile (ID).",
				Optional:    true,
				Default:     "",
			},
			"vmotion_network": {
				Type:        schema.TypeString,
				Description: "The vMotion network profile (ID).",
				Optional:    true,
				Default:     "",
			},
			"dvs": {
				Type:        schema.TypeString,
				Description: "The distributed switch used for L2 extension.",
				Required:    true,
			},
			"service": {
				Type:        schema.TypeList,
				Description: "The list of HCX services.",
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
		},
	}
}

func resourceComputeProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)

	res, err := hcx.GetVcInventory(client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Cluster info
	var cluster_id string
	var cluster_name string
	found := false
	for _, j := range res.Children[0].Children {
		if j.Name == cluster {
			cluster_id = j.Entity_id
			cluster_name = j.Name
			found = true
		}
	}
	if !found {
		return diag.FromErr(errors.New("cluster not found"))
	}

	// Get Datastore info
	datastore := d.Get("datastore").(string)
	datastore_from_api, err := hcx.GetVcDatastore(client, datastore, res.Entity_id, cluster_id)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get DVS info
	dvs := d.Get("dvs").(string)
	dvs_from_api, err := hcx.GetVcDvs(client, dvs, res.Entity_id, cluster_id)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Services from schema
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

	// Create network list with tags
	management_network := d.Get("management_network").(string)
	replication_network := d.Get("replication_network").(string)
	uplink_network := d.Get("uplink_network").(string)
	vmotion_network := d.Get("vmotion_network").(string)

	networks_list := []hcx.Network{}
	np, err := hcx.GetNetworkProfileById(client, management_network)
	if err != nil {
		return diag.FromErr(err)
	}
	management_network_name := np.Name
	management_network_id := np.ObjectId

	np, err = hcx.GetNetworkProfileById(client, replication_network)
	if err != nil {
		return diag.FromErr(err)
	}
	replication_network_name := np.Name
	replication_network_id := np.ObjectId

	np, err = hcx.GetNetworkProfileById(client, uplink_network)
	if err != nil {
		return diag.FromErr(err)
	}
	uplink_network_name := np.Name
	uplink_network_id := np.ObjectId

	np, err = hcx.GetNetworkProfileById(client, vmotion_network)
	if err != nil {
		return diag.FromErr(err)
	}
	vmotion_network_name := np.Name
	vmotion_network_id := np.ObjectId

	net_tmp := hcx.Network{
		Name: management_network_name,
		ID:   management_network_id,
		Tags: []string{"management"},
		Status: hcx.Status{
			State: "REALIZED",
		},
	}
	networks_list = append(networks_list, net_tmp)

	found = false
	index := 0
	for i, j := range networks_list {
		if j.Name == replication_network_name {
			found = true
			index = i
			break
		}
	}
	if found {
		networks_list[index].Tags = append(networks_list[index].Tags, "replication")
	} else {
		net_tmp := hcx.Network{
			Name: replication_network_name,
			ID:   replication_network_id,
			Tags: []string{"replication"},
			Status: hcx.Status{
				State: "REALIZED",
			},
		}
		networks_list = append(networks_list, net_tmp)
	}

	found = false
	index = 0
	for i, j := range networks_list {
		if j.Name == uplink_network_name {
			found = true
			index = i
			break
		}
	}
	if found {
		networks_list[index].Tags = append(networks_list[index].Tags, "uplink")
	} else {
		net_tmp := hcx.Network{
			Name: uplink_network_name,
			ID:   uplink_network_id,
			Tags: []string{"uplink"},
			Status: hcx.Status{
				State: "REALIZED",
			},
		}
		networks_list = append(networks_list, net_tmp)
	}

	found = false
	index = 0
	for i, j := range networks_list {
		if j.Name == vmotion_network_name {
			found = true
			index = i
			break
		}
	}
	if found {
		networks_list[index].Tags = append(networks_list[index].Tags, "vmotion")
	} else {
		net_tmp := hcx.Network{
			Name: vmotion_network_name,
			ID:   vmotion_network_id,
			Tags: []string{"vmotion"},
			Status: hcx.Status{
				State: "REALIZED",
			},
		}
		networks_list = append(networks_list, net_tmp)
	}

	body := hcx.InsertComputeProfileBody{
		Name:     name,
		Services: services_from_schema,
		Computes: []hcx.Compute{{
			CmpId:   res.Entity_id,
			CmpType: "VC",
			CmpName: res.Name,
			ID:      res.Children[0].Entity_id,
			Name:    res.Children[0].Name,
			Type:    res.Children[0].EntityType,
		}},
		DeploymentContainers: hcx.DeploymentContainer{
			Computes: []hcx.Compute{{
				CmpId:   res.Entity_id,
				CmpType: "VC",
				CmpName: res.Name,
				ID:      cluster_id,
				Name:    cluster_name,
				Type:    "ClusterComputeResource",
			},
			},
			Storage: []hcx.Storage{{
				CmpId:   res.Entity_id,
				CmpType: "VC",
				CmpName: res.Name,
				ID:      datastore_from_api.ID,
				Name:    datastore_from_api.Name,
				Type:    datastore_from_api.EntityType,
			}},
		},
		Networks: networks_list,
		Switches: []hcx.Switch{{
			CmpID:  res.Entity_id,
			MaxMTU: dvs_from_api.MaxMtu,
			ID:     dvs_from_api.ID,
			Name:   dvs_from_api.Name,
			Type:   dvs_from_api.Type,
		}},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	res2, err := hcx.InsertComputeProfile(client, body)

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

	d.SetId(res2.Data.ComputeProfileId)

	return resourceComputeProfileRead(ctx, d, m)

}

func resourceComputeProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceComputeProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceComputeProfileRead(ctx, d, m)
}

func resourceComputeProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	res, err := hcx.DeleteComputeProfile(client, d.Id())
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
