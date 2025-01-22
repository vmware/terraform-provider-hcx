// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Job_result struct {
	JobID                   string `json:"jobId"`
	Enterprise              string `json:"enterprise"`
	Organization            string `json:"organization"`
	Username                string `json:"username"`
	IsQueued                bool   `json:"isQueued"`
	IsCancelled             bool   `json:"isCancelled"`
	IsRolledBack            bool   `json:"isRolledBack"`
	CreateTimeEpoch         int64  `json:"createTimeEpoch"`
	AbsoluteExpireTimeEpoch int64  `json:"absoluteExpireTimeEpoch"`
	StartTime               int64  `json:"startTime"`
	EndTime                 int64  `json:"endTime"`
	PercentComplete         int    `json:"percentComplete"`
	IsDone                  bool   `json:"isDone"`
	DidFail                 bool   `json:"didFail"`
	TimeToExecute           int64  `json:"timeToExecute"`
}

type Task_result struct {
	InterconnectTaskId string `json:"interconnectTaskId"`
	Status             string `json:"status"`
}

type ResouceContainerListFilterCloud struct {
	Local  bool `json:"local"`
	Remote bool `json:"remote"`
}

type ResouceContainerListFilter struct {
	Cloud ResouceContainerListFilterCloud `json:"cloud"`
}

type PostResouceContainerListBody struct {
	Filter ResouceContainerListFilter `json:"filter"`
}

type PostResouceContainerListResult struct {
	Success   bool                               `json:"success"`
	Completed bool                               `json:"completed"`
	Time      int64                              `json:"time"`
	Data      PostResouceContainerListResultData `json:"data"`
}

type PostResouceContainerListResultData struct {
	Items []PostResouceContainerListResultDataItem `json:"items"`
}

type PostResouceContainerListResultDataItem struct {
	URL           string `json:"url"`
	Vcuuid        string `json:"vcuuid"`
	Version       string `json:"version"`
	BuildNumber   string `json:"buildNumber"`
	OsType        string `json:"osType"`
	Name          string `json:"name"`
	ResourceId    string `json:"resourceId"`
	ResourceType  string `json:"resourceType"`
	ResourceName  string `json:"resourceName"`
	VimId         string `json:"vimId"`
	VimServerUuid string `json:"vimServerUuid"`
}

type PostNetworkBackingBody struct {
	Filter PostNetworkBackingBodyFilter `json:"filter"`
}

type PostNetworkBackingBodyFilter struct {
	Cloud PostCloudListResultDataItem `json:"cloud"`
	//VCenterInstanceUuid string   `json:"vCenterInstanceUuid"`
	//ExcludeUsed         bool     `json:"excludeUsed"`
	//BackingTypes        []string `json:"backingTypes"`
}

type PostNetworkBackingResult struct {
	Data PostNetworkBackingResultData `json:"data"`
}

type PostNetworkBackingResultData struct {
	Items []Dvpg `json:"items"`
}

type Dvpg struct {
	EntityID   string `json:"entity_id"`
	Name       string `json:"name"`
	EntityType string `json:"entityType"`
}

type GetVcInventoryResult struct {
	Data GetVcInventoryResultData `json:"data"`
}

type GetVcInventoryResultData struct {
	Items []GetVcInventoryResultDataItem `json:"items"`
}

type GetVcInventoryResultDataItem struct {
	Vcenter_instanceId string                                 `json:"vcenter_instanceId"`
	Entity_id          string                                 `json:"entity_id"`
	Children           []GetVcInventoryResultDataItemChildren `json:"children"`
	Name               string                                 `json:"name"`
	EntityType         string                                 `json:"entityType"`
}

type GetVcInventoryResultDataItemChildren struct {
	Vcenter_instanceId string                                         `json:"vcenter_instanceId"`
	Entity_id          string                                         `json:"entity_id"`
	Children           []GetVcInventoryResultDataItemChildrenChildren `json:"children"`
	Name               string                                         `json:"name"`
	EntityType         string                                         `json:"entityType"`
}

type GetVcInventoryResultDataItemChildrenChildren struct {
	Vcenter_instanceId string `json:"vcenter_instanceId"`
	Entity_id          string `json:"entity_id"`
	Name               string `json:"name"`
	EntityType         string `json:"entityType"`
	// Datastores
}

type GetVcDatastoreResult struct {
	Success   bool                     `json:"success"`
	Completed bool                     `json:"completed"`
	Time      int64                    `json:"time"`
	Data      GetVcDatastoreResultData `json:"data"`
}

type GetVcDatastoreResultData struct {
	Items []GetVcDatastoreResultDataItem `json:"items"`
}

type GetVcDatastoreResultDataItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	EntityType string `json:"entity_type"`
}

type GetVcDatastoreBody struct {
	Filter GetVcDatastoreFilter `json:"filter"`
}

type GetVcDatastoreFilter struct {
	ComputeType       string   `json:"computeType"`
	VCenterInstanceID string   `json:"vcenter_instanceId"`
	ComputeIds        []string `json:"computeIds"`
}

type GetVcDvsResult struct {
	Success   bool               `json:"success"`
	Completed bool               `json:"completed"`
	Time      int64              `json:"time"`
	Data      GetVcDvsResultData `json:"data"`
}

type GetVcDvsResultData struct {
	Items []GetVcDvsResultDataItem `json:"items"`
}

type GetVcDvsResultDataItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	MaxMtu int    `json:"maxMtu"`
}

type GetVcDvsBody struct {
	Filter GetVcDvsFilter `json:"filter"`
}

type GetVcDvsFilter struct {
	ComputeType       string   `json:"computeType"`
	VCenterInstanceID string   `json:"vcenter_instanceId"`
	ComputeIds        []string `json:"computeIds"`
}

type PostCloudListFilter struct {
	Local  bool `json:"local"`
	Remote bool `json:"remote"`
}

type PostCloudListBody struct {
	Filter PostCloudListFilter `json:"filter"`
}

type PostCloudListResult struct {
	Success   bool                    `json:"success"`
	Completed bool                    `json:"completed"`
	Time      int64                   `json:"time"`
	Data      PostCloudListResultData `json:"data"`
}

type PostCloudListResultData struct {
	Items []PostCloudListResultDataItem `json:"items"`
}

type PostCloudListResultDataItem struct {
	EndpointId   string `json:"endpointId,omitempty"`
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	EndpointType string `json:"endpointType,omitempty"`
}

type GetApplianceBody struct {
	Filter GetApplianceBodyFilter `json:"filter"`
}

type GetApplianceBodyFilter struct {
	ApplianceType string `json:"applianceType"`
	EndpointId    string `json:"endpointId"`
	ServiceMeshId string `json:"serviceMeshId,omitempty"`
}

type GetApplianceResult struct {
	Items []GetApplianceResultItem `json:"items"`
}

type GetApplianceResultItem struct {
	ApplianceId           string `json:"applianceId"`
	ServiceMeshId         string `json:"serviceMeshId"`
	NetworkExtensionCount int    `json:"networkExtensionCount"`
}

// GetJobResult ...
func GetJobResult(c *Client, jobId string) (Job_result, error) {

	resp := Job_result{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/jobs/%s", c.HostURL, jobId), nil)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	return resp, nil
}

// GetTaskResult ...
func GetTaskResult(c *Client, taskId string) (Task_result, error) {

	resp := Task_result{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/interconnect/tasks/%s", c.HostURL, taskId), nil)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	return resp, nil
}

// GetLocalConatainer ...
func GetLocalContainer(c *Client) (PostResouceContainerListResultDataItem, error) {

	body := PostResouceContainerListBody{
		Filter: ResouceContainerListFilter{
			Cloud: ResouceContainerListFilterCloud{
				Local:  true,
				Remote: false,
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	resp := PostResouceContainerListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/resourcecontainer/list", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	return resp.Data.Items[0], nil
}

// GetLocalConatainer ...
func GetRemoteContainer(c *Client) (PostResouceContainerListResultDataItem, error) {

	body := PostResouceContainerListBody{
		Filter: ResouceContainerListFilter{
			Cloud: ResouceContainerListFilterCloud{
				Local:  false,
				Remote: true,
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	resp := PostResouceContainerListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/resourcecontainer/list", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostResouceContainerListResultDataItem{}, err
	}

	return resp.Data.Items[0], nil
}

// GetNetworkBacking ...
func GetNetworkBacking(c *Client, endpointid string, network string, network_type string) (Dvpg, error) {

	body := PostNetworkBackingBody{
		Filter: PostNetworkBackingBodyFilter{
			Cloud: PostCloudListResultDataItem{
				EndpointId: endpointid,
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return Dvpg{}, err
	}

	resp := PostNetworkBackingResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/networks", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return Dvpg{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return Dvpg{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return Dvpg{}, err
	}

	log.Printf("*************************************")
	log.Printf("networks list: %+v", resp)
	log.Printf("*************************************")

	for _, j := range resp.Data.Items {
		if j.Name == network && j.EntityType == network_type {
			return j, nil
		}
	}

	return Dvpg{}, errors.New("cannot find network info")
}

// GetVcInventory ...
func GetVcInventory(c *Client) (GetVcInventoryResultDataItem, error) {

	var jsonBody = []byte("{}")
	buf := bytes.NewBuffer(jsonBody)

	resp := GetVcInventoryResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/vc/list", c.HostURL), buf)
	if err != nil {
		fmt.Println(err)
		return GetVcInventoryResultDataItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetVcInventoryResultDataItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetVcInventoryResultDataItem{}, err
	}

	return resp.Data.Items[0], nil
}

// GetVcDatastore ...
func GetVcDatastore(c *Client, datastore_name string, vcuuid string, cluster string) (GetVcDatastoreResultDataItem, error) {

	body := GetVcDatastoreBody{
		Filter: GetVcDatastoreFilter{
			VCenterInstanceID: vcuuid,
			ComputeType:       "ClusterComputeResource",
			ComputeIds:        []string{cluster},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return GetVcDatastoreResultDataItem{}, err
	}

	resp := GetVcDatastoreResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/vc/datastores/query", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return GetVcDatastoreResultDataItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetVcDatastoreResultDataItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetVcDatastoreResultDataItem{}, err
	}

	for _, j := range resp.Data.Items {
		if j.Name == datastore_name {
			return j, nil
		}
	}

	return GetVcDatastoreResultDataItem{}, errors.New("cannot find datastore")
}

// GetVcDvs ...
func GetVcDvs(c *Client, dvs_name string, vcuuid string, cluster string) (GetVcDvsResultDataItem, error) {

	body := GetVcDvsBody{
		Filter: GetVcDvsFilter{
			VCenterInstanceID: vcuuid,
			ComputeType:       "ClusterComputeResource",
			ComputeIds:        []string{cluster},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return GetVcDvsResultDataItem{}, err
	}

	resp := GetVcDvsResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/vc/dvs/query", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return GetVcDvsResultDataItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetVcDvsResultDataItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetVcDvsResultDataItem{}, err
	}

	for _, j := range resp.Data.Items {
		if j.Name == dvs_name {
			return j, nil
		}
	}

	return GetVcDvsResultDataItem{}, errors.New("cannot find datastore")
}

// GetRemoteCloudList ...
func GetRemoteCloudList(c *Client) (PostCloudListResult, error) {

	body := PostCloudListBody{
		Filter: PostCloudListFilter{
			Remote: true,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	resp := PostCloudListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/cloud/list", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	return resp, nil
}

// GetRemoteCloudList ...
func GetLocalCloudList(c *Client) (PostCloudListResult, error) {

	body := PostCloudListBody{
		Filter: PostCloudListFilter{
			Local: true,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	resp := PostCloudListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/cloud/list", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	return resp, nil
}

// GetRemoteCloudList ...
func GetAppliance(c *Client, endpointId string, service_mesh_id string) (GetApplianceResultItem, error) {

	body := GetApplianceBody{
		Filter: GetApplianceBodyFilter{
			ApplianceType: "HCX-NET-EXT",
			EndpointId:    endpointId,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return GetApplianceResultItem{}, err
	}

	resp := GetApplianceResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/appliances/query", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return GetApplianceResultItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetApplianceResultItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetApplianceResultItem{}, err
	}

	for _, j := range resp.Items {
		if j.ServiceMeshId == service_mesh_id && j.NetworkExtensionCount < 9 {
			return j, nil
		}
	}

	return resp.Items[0], nil
}

// GetRemoteCloudList ...
func GetAppliances(c *Client, endpointId string, service_mesh_id string) ([]GetApplianceResultItem, error) {

	body := GetApplianceBody{
		Filter: GetApplianceBodyFilter{
			ApplianceType: "HCX-NET-EXT",
			EndpointId:    endpointId,
			ServiceMeshId: service_mesh_id,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return []GetApplianceResultItem{}, err
	}

	resp := GetApplianceResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/appliances/query", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return []GetApplianceResultItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return []GetApplianceResultItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return []GetApplianceResultItem{}, err
	}

	return resp.Items, nil
}
