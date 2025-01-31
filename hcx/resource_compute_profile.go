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

// resourceComputeProfile defines the resource schema for managing compute profile configuration.
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

// resourceComputeProfileCreate creates the compute profile configuration.
func resourceComputeProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)

	res, err := GetVcInventory(client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Cluster info
	var clusterID string
	var clusterName string
	found := false
	for _, j := range res.Children[0].Children {
		if j.Name == cluster {
			clusterID = j.EntityID
			clusterName = j.Name
			found = true
		}
	}
	if !found {
		return diag.FromErr(errors.New("cluster not found"))
	}

	// Get Datastore info
	datastore := d.Get("datastore").(string)
	datastoreFromAPI, err := GetVcDatastore(client, datastore, res.EntityID, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get DVS info
	dvs := d.Get("dvs").(string)
	dvsFromAPI, err := GetVcDvs(client, dvs, res.EntityID, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Services from schema
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

	// Create network list with tags
	managementNetwork := d.Get("management_network").(string)
	replicationNetwork := d.Get("replication_network").(string)
	uplinkNetwork := d.Get("uplink_network").(string)
	vmotionNetwork := d.Get("vmotion_network").(string)

	networksList := []Network{}
	np, err := GetNetworkProfileByID(client, managementNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	managementNetworkName := np.Name
	managementNetworkID := np.ObjectID

	np, err = GetNetworkProfileByID(client, replicationNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	replicationNetworkName := np.Name
	replicationNetworkID := np.ObjectID

	np, err = GetNetworkProfileByID(client, uplinkNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	uplinkNetworkName := np.Name
	uplinkNetworkID := np.ObjectID

	np, err = GetNetworkProfileByID(client, vmotionNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	vmotionNetworkName := np.Name
	vmotionNetworkID := np.ObjectID

	netTmp := Network{
		Name: managementNetworkName,
		ID:   managementNetworkID,
		Tags: []string{"management"},
		Status: Status{
			State: constants.RealizedStatus,
		},
	}
	networksList = append(networksList, netTmp)

	found = false
	index := 0
	for i, j := range networksList {
		if j.Name == replicationNetworkName {
			found = true
			index = i
			break
		}
	}
	if found {
		networksList[index].Tags = append(networksList[index].Tags, "replication")
	} else {
		netTmp := Network{
			Name: replicationNetworkName,
			ID:   replicationNetworkID,
			Tags: []string{"replication"},
			Status: Status{
				State: constants.RealizedStatus,
			},
		}
		networksList = append(networksList, netTmp)
	}

	found = false
	index = 0
	for i, j := range networksList {
		if j.Name == uplinkNetworkName {
			found = true
			index = i
			break
		}
	}
	if found {
		networksList[index].Tags = append(networksList[index].Tags, "uplink")
	} else {
		netTmp := Network{
			Name: uplinkNetworkName,
			ID:   uplinkNetworkID,
			Tags: []string{"uplink"},
			Status: Status{
				State: constants.RealizedStatus,
			},
		}
		networksList = append(networksList, netTmp)
	}

	found = false
	index = 0
	for i, j := range networksList {
		if j.Name == vmotionNetworkName {
			found = true
			index = i
			break
		}
	}
	if found {
		networksList[index].Tags = append(networksList[index].Tags, "vmotion")
	} else {
		netTmp := Network{
			Name: vmotionNetworkName,
			ID:   vmotionNetworkID,
			Tags: []string{"vmotion"},
			Status: Status{
				State: constants.RealizedStatus,
			},
		}
		networksList = append(networksList, netTmp)
	}

	body := InsertComputeProfileBody{
		Name:     name,
		Services: servicesFromSchema,
		Computes: []Compute{{
			ComputeID:   res.EntityID,
			ComputeType: constants.DefaultComputeType,
			ComputeName: res.Name,
			ID:          res.Children[0].EntityID,
			Name:        res.Children[0].Name,
			Type:        res.Children[0].EntityType,
		}},
		DeploymentContainers: DeploymentContainer{
			Computes: []Compute{{
				ComputeID:   res.EntityID,
				ComputeType: constants.DefaultComputeType,
				ComputeName: res.Name,
				ID:          clusterID,
				Name:        clusterName,
				Type:        "ClusterComputeResource",
			},
			},
			Storage: []Storage{{
				ComputeID:   res.EntityID,
				ComputeType: constants.DefaultComputeType,
				ComputeName: res.Name,
				ID:          datastoreFromAPI.ID,
				Name:        datastoreFromAPI.Name,
				Type:        datastoreFromAPI.EntityType,
			}},
		},
		Networks: networksList,
		Switches: []Switch{{
			ComputeID: res.EntityID,
			MaxMTU:    dvsFromAPI.MaxMTU,
			ID:        dvsFromAPI.ID,
			Name:      dvsFromAPI.Name,
			Type:      dvsFromAPI.Type,
		}},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return diag.FromErr(err)
	}

	res2, err := InsertComputeProfile(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task completion
	for {
		jr, err := GetTaskResult(client, res2.Data.InterconnectTaskID)
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

	d.SetId(res2.Data.ComputeProfileID)

	return resourceComputeProfileRead(ctx, d, m)

}

// resourceComputeProfileRead retrieves the compute profile configuration.
func resourceComputeProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

// resourceComputeProfileUpdate updates the compute profile configuration.
func resourceComputeProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceComputeProfileRead(ctx, d, m)
}

// resourceComputeProfileDelete removes the compute profile configuration and clears the state of the resource in the
// schema.
func resourceComputeProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	res, err := DeleteComputeProfile(client, d.Id())
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
