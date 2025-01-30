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

type JobResult struct {
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

type TaskResult struct {
	InterconnectTaskID string `json:"interconnectTaskId"`
	Status             string `json:"status"`
}

type ResourceContainerListFilterCloud struct {
	Local  bool `json:"local"`
	Remote bool `json:"remote"`
}

type ResourceContainerListFilter struct {
	Cloud ResourceContainerListFilterCloud `json:"cloud"`
}

type PostResourceContainerListBody struct {
	Filter ResourceContainerListFilter `json:"filter"`
}

type PostResourceContainerListResult struct {
	Success   bool                                `json:"success"`
	Completed bool                                `json:"completed"`
	Time      int64                               `json:"time"`
	Data      PostResourceContainerListResultData `json:"data"`
}

type PostResourceContainerListResultData struct {
	Items []PostResourceContainerListResultDataItem `json:"items"`
}

type PostResourceContainerListResultDataItem struct {
	URL           string `json:"url"`
	VcUUID        string `json:"vcuuid"`
	Version       string `json:"version"`
	BuildNumber   string `json:"buildNumber"`
	OsType        string `json:"osType"`
	Name          string `json:"name"`
	ResourceID    string `json:"resourceId"`
	ResourceType  string `json:"resourceType"`
	ResourceName  string `json:"resourceName"`
	VimID         string `json:"vimId"`
	VimServerUUID string `json:"vimServerUuid"`
}

type PostNetworkBackingBody struct {
	Filter PostNetworkBackingBodyFilter `json:"filter"`
}

type PostNetworkBackingBodyFilter struct {
	Cloud PostCloudListResultDataItem `json:"cloud"`
	//VCenterInstanceUUID string   `json:"vCenterInstanceUuid"`
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
	VCenterInstanceID string                                 `json:"vcenter_instanceId"`
	EntityID          string                                 `json:"entity_id"`
	Children          []GetVcInventoryResultDataItemChildren `json:"children"`
	Name              string                                 `json:"name"`
	EntityType        string                                 `json:"entityType"`
}

type GetVcInventoryResultDataItemChildren struct {
	VCenterInstanceID string                                         `json:"vcenter_instanceId"`
	EntityID          string                                         `json:"entity_id"`
	Children          []GetVcInventoryResultDataItemChildrenChildren `json:"children"`
	Name              string                                         `json:"name"`
	EntityType        string                                         `json:"entityType"`
}

type GetVcInventoryResultDataItemChildrenChildren struct {
	VCenterInstanceID string `json:"vcenter_instanceId"`
	EntityID          string `json:"entity_id"`
	Name              string `json:"name"`
	EntityType        string `json:"entityType"`
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
	ComputeIDs        []string `json:"computeIds"`
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
	MaxMTU int    `json:"maxMtu"`
}

type GetVcDvsBody struct {
	Filter GetVcDvsFilter `json:"filter"`
}

type GetVcDvsFilter struct {
	ComputeType       string   `json:"computeType"`
	VCenterInstanceID string   `json:"vcenter_instanceId"`
	ComputeIDs        []string `json:"computeIds"`
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
	EndpointID   string `json:"endpointId,omitempty"`
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	EndpointType string `json:"endpointType,omitempty"`
}

type GetApplianceBody struct {
	Filter GetApplianceBodyFilter `json:"filter"`
}

type GetApplianceBodyFilter struct {
	ApplianceType string `json:"applianceType"`
	EndpointID    string `json:"endpointId"`
	ServiceMeshID string `json:"serviceMeshId,omitempty"`
}

type GetApplianceResult struct {
	Items []GetApplianceResultItem `json:"items"`
}

type GetApplianceResultItem struct {
	ApplianceID           string `json:"applianceId"`
	ServiceMeshID         string `json:"serviceMeshId"`
	NetworkExtensionCount int    `json:"networkExtensionCount"`
}

// GetJobResult ...
func GetJobResult(c *Client, jobID string) (JobResult, error) {

	resp := JobResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/jobs/%s", c.HostURL, jobID), nil)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	return resp, nil
}

// GetTaskResult ...
func GetTaskResult(c *Client, taskID string) (TaskResult, error) {

	resp := TaskResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/interconnect/tasks/%s", c.HostURL, taskID), nil)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	return resp, nil
}

// GetLocalConatainer ...
func GetLocalContainer(c *Client) (PostResourceContainerListResultDataItem, error) {

	body := PostResourceContainerListBody{
		Filter: ResourceContainerListFilter{
			Cloud: ResourceContainerListFilterCloud{
				Local:  true,
				Remote: false,
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	resp := PostResourceContainerListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/resourcecontainer/list", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	return resp.Data.Items[0], nil
}

// GetLocalConatainer ...
func GetRemoteContainer(c *Client) (PostResourceContainerListResultDataItem, error) {

	body := PostResourceContainerListBody{
		Filter: ResourceContainerListFilter{
			Cloud: ResourceContainerListFilterCloud{
				Local:  false,
				Remote: true,
			},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	resp := PostResourceContainerListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/resourcecontainer/list", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostResourceContainerListResultDataItem{}, err
	}

	return resp.Data.Items[0], nil
}

// GetNetworkBacking ...
func GetNetworkBacking(c *Client, endpointID string, network string, networkType string) (Dvpg, error) {

	body := PostNetworkBackingBody{
		Filter: PostNetworkBackingBodyFilter{
			Cloud: PostCloudListResultDataItem{
				EndpointID: endpointID,
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return Dvpg{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return Dvpg{}, err
	}

	log.Printf("*************************************")
	log.Printf("networks list: %+v", resp)
	log.Printf("*************************************")

	for _, j := range resp.Data.Items {
		if j.Name == network && j.EntityType == networkType {
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetVcInventoryResultDataItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetVcInventoryResultDataItem{}, err
	}

	return resp.Data.Items[0], nil
}

// GetVcDatastore ...
func GetVcDatastore(c *Client, datastoreName string, vcuuid string, cluster string) (GetVcDatastoreResultDataItem, error) {

	body := GetVcDatastoreBody{
		Filter: GetVcDatastoreFilter{
			VCenterInstanceID: vcuuid,
			ComputeType:       "ClusterComputeResource",
			ComputeIDs:        []string{cluster},
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetVcDatastoreResultDataItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetVcDatastoreResultDataItem{}, err
	}

	for _, j := range resp.Data.Items {
		if j.Name == datastoreName {
			return j, nil
		}
	}

	return GetVcDatastoreResultDataItem{}, errors.New("cannot find datastore")
}

// GetVcDvs ...
func GetVcDvs(c *Client, dvsName string, vcuuid string, cluster string) (GetVcDvsResultDataItem, error) {

	body := GetVcDvsBody{
		Filter: GetVcDvsFilter{
			VCenterInstanceID: vcuuid,
			ComputeType:       "ClusterComputeResource",
			ComputeIDs:        []string{cluster},
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetVcDvsResultDataItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetVcDvsResultDataItem{}, err
	}

	for _, j := range resp.Data.Items {
		if j.Name == dvsName {
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	// Parse response body.
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return PostCloudListResult{}, err
	}

	return resp, nil
}

// GetRemoteCloudList ...
func GetAppliance(c *Client, endpointID string, serviceMeshID string) (GetApplianceResultItem, error) {

	body := GetApplianceBody{
		Filter: GetApplianceBodyFilter{
			ApplianceType: "HCX-NET-EXT",
			EndpointID:    endpointID,
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetApplianceResultItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetApplianceResultItem{}, err
	}

	for _, j := range resp.Items {
		if j.ServiceMeshID == serviceMeshID && j.NetworkExtensionCount < 9 {
			return j, nil
		}
	}

	return resp.Items[0], nil
}

// GetRemoteCloudList ...
func GetAppliances(c *Client, endpointID string, serviceMeshID string) ([]GetApplianceResultItem, error) {

	body := GetApplianceBody{
		Filter: GetApplianceBodyFilter{
			ApplianceType: "HCX-NET-EXT",
			EndpointID:    endpointID,
			ServiceMeshID: serviceMeshID,
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

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return []GetApplianceResultItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return []GetApplianceResultItem{}, err
	}

	return resp.Items, nil
}
