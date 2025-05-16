// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// JobResult represents the result of a job execution, including its status, timing, and completion details.
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

// TaskResult represents the result of a run.
type TaskResult struct {
	InterconnectTaskID string `json:"interconnectTaskId"`
	Status             string `json:"status"`
}

// ResourceContainerListFilterCloud defines a filter structure for categorizing resource containers as local or remote.
type ResourceContainerListFilterCloud struct {
	Local  bool `json:"local"`
	Remote bool `json:"remote"`
}

// ResourceContainerListFilter defines a structure for filtering resource containers based on cloud characteristics.
type ResourceContainerListFilter struct {
	Cloud ResourceContainerListFilterCloud `json:"cloud"`
}

// PostResourceContainerListBody defines the request body for retrieving a filtered list of resource containers.
type PostResourceContainerListBody struct {
	Filter ResourceContainerListFilter `json:"filter"`
}

// PostResourceContainerListResult represents the result of a resource container list request, including success status.
type PostResourceContainerListResult struct {
	Success   bool                                `json:"success"`
	Completed bool                                `json:"completed"`
	Time      int64                               `json:"time"`
	Data      PostResourceContainerListResultData `json:"data"`
}

// PostResourceContainerListResultData represents the data structure for a container list result containing resource
// items.
type PostResourceContainerListResultData struct {
	Items []PostResourceContainerListResultDataItem `json:"items"`
}

// PostResourceContainerListResultDataItem defines a structure representing a single resource container in the list.
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

// PostNetworkBackingResult represents the structure for the response of a network backing operation, containing data.
type PostNetworkBackingResult struct {
	Data PostNetworkBackingResultData `json:"data"`
}

// PostNetworkBackingResultData represents the response containing a list of distributed port groups.
type PostNetworkBackingResultData struct {
	Items []Dvpg `json:"items"`
}

// Dvpg represents a distributed port group with associated metadata.
type Dvpg struct {
	EntityID   string `json:"entity_id"`
	Name       string `json:"name"`
	EntityType string `json:"entityType"`
}

// GetVcInventoryResult represents the structure for the vCenter inventory result containing data about inventory items.
type GetVcInventoryResult struct {
	Data GetVcInventoryResultData `json:"data"`
}

// GetVcInventoryResultData represents the structure containing a list of vCenter instance inventory result data items.
type GetVcInventoryResultData struct {
	Items []GetVcInventoryResultDataItem `json:"items"`
}

// GetVcInventoryResultDataItem represents an inventory item in a vCenter instance including its children and metadata.
type GetVcInventoryResultDataItem struct {
	VCenterInstanceID string                                 `json:"vcenter_instanceId"`
	EntityID          string                                 `json:"entity_id"`
	Children          []GetVcInventoryResultDataItemChildren `json:"children"`
	Name              string                                 `json:"name"`
	EntityType        string                                 `json:"entityType"`
}

// GetVcInventoryResultDataItemChildren represents an inventory item in vCenter instance with nested children and
// metadata.
type GetVcInventoryResultDataItemChildren struct {
	VCenterInstanceID string                                         `json:"vcenter_instanceId"`
	EntityID          string                                         `json:"entity_id"`
	Children          []GetVcInventoryResultDataItemChildrenChildren `json:"children"`
	Name              string                                         `json:"name"`
	EntityType        string                                         `json:"entityType"`
}

// GetVcInventoryResultDataItemChildrenChildren represents nested child items in a vCenter instance inventory structure.
type GetVcInventoryResultDataItemChildrenChildren struct {
	VCenterInstanceID string `json:"vcenter_instanceId"`
	EntityID          string `json:"entity_id"`
	Name              string `json:"name"`
	EntityType        string `json:"entityType"`
	// Datastores
}

// GetVcDatastoreResult represents the result of a query for vCenter instance datastore information.
type GetVcDatastoreResult struct {
	Success   bool                     `json:"success"`
	Completed bool                     `json:"completed"`
	Time      int64                    `json:"time"`
	Data      GetVcDatastoreResultData `json:"data"`
}

// GetVcDatastoreResultData represents a collection of datastore result items retrieved from a vCenter instance.
type GetVcDatastoreResultData struct {
	Items []GetVcDatastoreResultDataItem `json:"items"`
}

// GetVcDatastoreResultDataItem represents a single datastore entity retrieved from the vCenter instance.
type GetVcDatastoreResultDataItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	EntityType string `json:"entity_type"`
}

// GetVcDatastoreBody represents the body of the request for querying vCenter instance datastores using specified
// filters.
type GetVcDatastoreBody struct {
	Filter GetVcDatastoreFilter `json:"filter"`
}

// GetVcDatastoreFilter defines a filter for querying vCenter instance datastoress.
type GetVcDatastoreFilter struct {
	ComputeType       string   `json:"computeType"`
	VCenterInstanceID string   `json:"vcenter_instanceId"`
	ComputeIDs        []string `json:"computeIds"`
}

// GetVcDvsResult represents the result of querying a distributed switch in a vCenter instance.
type GetVcDvsResult struct {
	Success   bool               `json:"success"`
	Completed bool               `json:"completed"`
	Time      int64              `json:"time"`
	Data      GetVcDvsResultData `json:"data"`
}

// GetVcDvsResultData represents a collection of distributed switch data retrieved from a vCenter instance.
type GetVcDvsResultData struct {
	Items []GetVcDvsResultDataItem `json:"items"`
}

// GetVcDvsResultDataItem represents an individual distributed virtual switch retrieved from a vCenter instance.
type GetVcDvsResultDataItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	MaxMTU int    `json:"maxMtu"`
}

// GetVcDvsBody represents the request body containing parameters to query a distributed switch.
type GetVcDvsBody struct {
	Filter GetVcDvsFilter `json:"filter"`
}

// GetVcDvsFilter represents the filter criteria for querying distributed switches in a vCenter instance.
type GetVcDvsFilter struct {
	ComputeType       string   `json:"computeType"`
	VCenterInstanceID string   `json:"vcenter_instanceId"`
	ComputeIDs        []string `json:"computeIds"`
}

// PostCloudListFilter is a struct used to filter cloud list requests based on "local" or "remote" cloud availability.
type PostCloudListFilter struct {
	Local  bool `json:"local"`
	Remote bool `json:"remote"`
}

// PostCloudListBody represents the body of a request to retrieve a filtered list of clouds.
type PostCloudListBody struct {
	Filter PostCloudListFilter `json:"filter"`
}

// PostCloudListResult contains the structure for response data from a cloud list request.
type PostCloudListResult struct {
	Success   bool                    `json:"success"`
	Completed bool                    `json:"completed"`
	Time      int64                   `json:"time"`
	Data      PostCloudListResultData `json:"data"`
}

// PostCloudListResultData represents a collection of cloud endpoint data items in a response.
type PostCloudListResultData struct {
	Items []PostCloudListResultDataItem `json:"items"`
}

// PostCloudListResultDataItem represents a single cloud endpoint in the response data.
type PostCloudListResultDataItem struct {
	EndpointID   string `json:"endpointId,omitempty"`
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	EndpointType string `json:"endpointType,omitempty"`
}

// GetApplianceBody represents the request body structure for querying appliances based on filters.
type GetApplianceBody struct {
	Filter GetApplianceBodyFilter `json:"filter"`
}

// GetApplianceBodyFilter defines the structure for filtering appliance queries.
type GetApplianceBodyFilter struct {
	ApplianceType string `json:"applianceType"`
	EndpointID    string `json:"endpointId"`
	ServiceMeshID string `json:"serviceMeshId,omitempty"`
}

// GetApplianceResult represents the result of querying appliances.
type GetApplianceResult struct {
	Items []GetApplianceResultItem `json:"items"`
}

// GetApplianceResultItem represents an appliance query result.
type GetApplianceResultItem struct {
	ApplianceID           string `json:"applianceId"`
	ServiceMeshID         string `json:"serviceMeshId"`
	NetworkExtensionCount int    `json:"networkExtensionCount"`
}

// GetJobResult sends a request to retrieve the result of a job identified by the provided jobID, returning a JobResult object.
func GetJobResult(c *Client, jobID string) (JobResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting job result", map[string]interface{}{
		"jobID": jobID,
	})

	resp := JobResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/jobs/%s", c.HostURL, jobID), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create job result GET request", map[string]interface{}{
			"error": err.Error(),
			"jobID": jobID,
		})
		return resp, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send job result GET request", map[string]interface{}{
			"error": err.Error(),
			"jobID": jobID,
		})
		return resp, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse job result response", map[string]interface{}{
			"error": err.Error(),
			"jobID": jobID,
		})
		return resp, fmt.Errorf("failed to parse job result response: %w", err)
	}

	tflog.Debug(ctx, "Job result retrieved successfully", map[string]interface{}{
		"jobID":           jobID,
		"percentComplete": resp.PercentComplete,
		"isDone":          resp.IsDone,
		"didFail":         resp.DidFail,
	})
	return resp, nil
}

// GetTaskResult sends a request to retrieve the result of a task identified by the provided taskID.
func GetTaskResult(c *Client, taskID string) (TaskResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting task result", map[string]interface{}{
		"taskID": taskID,
	})

	resp := TaskResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/interconnect/tasks/%s", c.HostURL, taskID), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create task result GET request", map[string]interface{}{
			"error":  err.Error(),
			"taskID": taskID,
		})
		return resp, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send task result GET request", map[string]interface{}{
			"error":  err.Error(),
			"taskID": taskID,
		})
		return resp, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse task result response", map[string]interface{}{
			"error":  err.Error(),
			"taskID": taskID,
		})
		return resp, fmt.Errorf("failed to parse task result response: %w", err)
	}

	tflog.Debug(ctx, "Task result retrieved successfully", map[string]interface{}{
		"taskID": taskID,
		"status": resp.Status,
	})
	return resp, nil
}

// GetLocalContainer sends a request to retrieve the local resource container list and returns the first item as a
// PostResourceContainerListResultDataItem.
func GetLocalContainer(c *Client) (PostResourceContainerListResultDataItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting local resource container")

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
		tflog.Error(ctx, "Failed to encode local container request body", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := PostResourceContainerListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/resourcecontainer/list", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create local container POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send local container HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse local container HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "Local container retrieved successfully", map[string]interface{}{
		"container_name": resp.Data.Items[0].Name,
		"resource_id":    resp.Data.Items[0].ResourceID,
	})
	return resp.Data.Items[0], nil
}

// GetRemoteContainer sends a request to retrieve the remote resource container list and returns the first item as a
// PostResourceContainerListResultDataItem.
func GetRemoteContainer(c *Client) (PostResourceContainerListResultDataItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting remote resource container")

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
		tflog.Error(ctx, "Failed to encode remote container request body", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := PostResourceContainerListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/resourcecontainer/list", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create remote container POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send remote container HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse remote container HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return PostResourceContainerListResultDataItem{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "Remote container retrieved successfully", map[string]interface{}{
		"container_name": resp.Data.Items[0].Name,
		"resource_id":    resp.Data.Items[0].ResourceID,
	})
	return resp.Data.Items[0], nil
}

// GetNetworkBacking sends a request to retrieve a network's backing information.
func GetNetworkBacking(c *Client, endpointID, network, networkType string) (Dvpg, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting network backing information", map[string]interface{}{
		"endpointID":  endpointID,
		"network":     network,
		"networkType": networkType,
	})

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
		tflog.Error(ctx, "Failed to encode network backing request body", map[string]interface{}{
			"error":      err.Error(),
			"endpointID": endpointID,
		})
		return Dvpg{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := PostNetworkBackingResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/networks", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create network backing POST request", map[string]interface{}{
			"error":      err.Error(),
			"endpointID": endpointID,
		})
		return Dvpg{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send network backing POST request", map[string]interface{}{
			"error":      err.Error(),
			"endpointID": endpointID,
		})
		return Dvpg{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse network backing HTTP response", map[string]interface{}{
			"error":      err.Error(),
			"endpointID": endpointID,
		})
		return Dvpg{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	for _, j := range resp.Data.Items {
		if j.Name == network && j.EntityType == networkType {
			tflog.Debug(ctx, "Found matching network backing", map[string]interface{}{
				"network":     network,
				"networkType": networkType,
				"entityID":    j.EntityID,
			})
			return j, nil
		}
	}

	tflog.Error(ctx, "Failed to find matching network info", map[string]interface{}{
		"network":     network,
		"networkType": networkType,
	})
	return Dvpg{}, fmt.Errorf("failed to find matching network info in GetNetworkBacking")
}

// GetVcInventory sends a request to retrieve the vCenter resource inventory.
func GetVcInventory(c *Client) (GetVcInventoryResultDataItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting vCenter inventory")

	var jsonBody = []byte("{}")
	buf := bytes.NewBuffer(jsonBody)

	resp := GetVcInventoryResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/vc/list", c.HostURL), buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create vCenter inventory POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return GetVcInventoryResultDataItem{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send vCenter inventory POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return GetVcInventoryResultDataItem{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse vCenter inventory HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return GetVcInventoryResultDataItem{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "vCenter inventory retrieved successfully", map[string]interface{}{
		"vcenter_instance_id": resp.Data.Items[0].VCenterInstanceID,
		"name":                resp.Data.Items[0].Name,
	})
	return resp.Data.Items[0], nil
}

// GetVcDatastore sends a request to query a vCenter datastore.
func GetVcDatastore(c *Client, datastoreName, vcuuid, cluster string) (GetVcDatastoreResultDataItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting vCenter datastore", map[string]interface{}{
		"datastoreName": datastoreName,
		"vcuuid":        vcuuid,
		"cluster":       cluster,
	})

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
		tflog.Error(ctx, "Failed to encode vCenter datastore request body", map[string]interface{}{
			"error":         err.Error(),
			"datastoreName": datastoreName,
			"vcuuid":        vcuuid,
		})
		return GetVcDatastoreResultDataItem{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := GetVcDatastoreResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/vc/datastores/query", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create vCenter datastore POST request", map[string]interface{}{
			"error":         err.Error(),
			"datastoreName": datastoreName,
		})
		return GetVcDatastoreResultDataItem{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send vCenter datastore POST request", map[string]interface{}{
			"error":         err.Error(),
			"datastoreName": datastoreName,
		})
		return GetVcDatastoreResultDataItem{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse vCenter datastore HTTP response", map[string]interface{}{
			"error":         err.Error(),
			"datastoreName": datastoreName,
		})
		return GetVcDatastoreResultDataItem{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	for _, j := range resp.Data.Items {
		if j.Name == datastoreName {
			tflog.Debug(ctx, "Found matching datastore", map[string]interface{}{
				"datastoreName": datastoreName,
				"id":            j.ID,
			})
			return j, nil
		}
	}

	tflog.Error(ctx, "Datastore not found", map[string]interface{}{
		"datastoreName": datastoreName,
	})
	return GetVcDatastoreResultDataItem{}, fmt.Errorf("datastore not found in GetVcDatastore")
}

// GetVcDvs sends a request to query a distributed switch.
func GetVcDvs(c *Client, dvsName, vcuuid, cluster string) (GetVcDvsResultDataItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting vCenter distributed switch", map[string]interface{}{
		"dvsName": dvsName,
		"vcuuid":  vcuuid,
		"cluster": cluster,
	})

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
		tflog.Error(ctx, "Failed to encode vCenter distributed switch request body", map[string]interface{}{
			"error":   err.Error(),
			"dvsName": dvsName,
			"vcuuid":  vcuuid,
		})
		return GetVcDvsResultDataItem{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := GetVcDvsResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/vc/dvs/query", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create vCenter distributed switch POST request", map[string]interface{}{
			"error":   err.Error(),
			"dvsName": dvsName,
		})
		return GetVcDvsResultDataItem{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send vCenter distributed switch POST request", map[string]interface{}{
			"error":   err.Error(),
			"dvsName": dvsName,
		})
		return GetVcDvsResultDataItem{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse vCenter distributed switch HTTP response", map[string]interface{}{
			"error":   err.Error(),
			"dvsName": dvsName,
		})
		return GetVcDvsResultDataItem{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	for _, j := range resp.Data.Items {
		if j.Name == dvsName {
			tflog.Debug(ctx, "Found matching distributed switch", map[string]interface{}{
				"dvsName": dvsName,
				"id":      j.ID,
				"maxMTU":  j.MaxMTU,
			})
			return j, nil
		}
	}

	tflog.Error(ctx, "Distributed switch not found", map[string]interface{}{
		"dvsName": dvsName,
	})
	return GetVcDvsResultDataItem{}, fmt.Errorf("distributed switch not found in GetVcDvs")
}

// GetRemoteCloudList sends a request to retrieve a list of remote clouds and returns the resulting PostCloudListResult.
func GetRemoteCloudList(c *Client) (PostCloudListResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting remote cloud list")

	body := PostCloudListBody{
		Filter: PostCloudListFilter{
			Remote: true,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode remote cloud list request body", map[string]interface{}{
			"error": err.Error(),
		})
		return PostCloudListResult{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := PostCloudListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/cloud/list", c.HostURL), &buf)
	if err != nil {
		return PostCloudListResult{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return PostCloudListResult{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return PostCloudListResult{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}

// GetLocalCloudList sends a request to retrieve a list of local clouds and returns the resulting PostCloudListResult.
func GetLocalCloudList(c *Client) (PostCloudListResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting local cloud list")

	body := PostCloudListBody{
		Filter: PostCloudListFilter{
			Local: true,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode local cloud list request body", map[string]interface{}{
			"error": err.Error(),
		})
		return PostCloudListResult{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := PostCloudListResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/service/inventory/cloud/list", c.HostURL), &buf)
	if err != nil {
		return PostCloudListResult{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return PostCloudListResult{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return PostCloudListResult{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}

// GetAppliance sends a request to query appliances based on the given endpointID and serviceMeshID.
// It returns a matching GetApplianceResultItem (with a network extension count less than 9) or an error.
func GetAppliance(c *Client, endpointID string, serviceMeshID string) (GetApplianceResultItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting appliance information", map[string]interface{}{
		"endpointID":    endpointID,
		"serviceMeshID": serviceMeshID,
	})

	body := GetApplianceBody{
		Filter: GetApplianceBodyFilter{
			ApplianceType: "HCX-NET-EXT",
			EndpointID:    endpointID,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode appliance request body", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return GetApplianceResultItem{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	resp := GetApplianceResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/appliances/query", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create appliance POST request", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return GetApplianceResultItem{}, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send appliance POST request", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return GetApplianceResultItem{}, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse appliance HTTP response", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return GetApplianceResultItem{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	for _, j := range resp.Items {
		if j.ServiceMeshID == serviceMeshID && j.NetworkExtensionCount < 9 {
			tflog.Debug(ctx, "Found matching appliance", map[string]interface{}{
				"endpointID":            endpointID,
				"serviceMeshID":         serviceMeshID,
				"networkExtensionCount": j.NetworkExtensionCount,
			})
			return j, nil
		}
	}

	tflog.Debug(ctx, "Returning first appliance in list", map[string]interface{}{
		"endpointID":    endpointID,
		"serviceMeshID": serviceMeshID,
	})
	return resp.Items[0], nil
}

// GetAppliances sends a request to retrieve all appliances matching the given endpointID and serviceMeshID.
// Returns a slice of GetApplianceResultItem objects or an error if the request fails.
func GetAppliances(c *Client, endpointID string, serviceMeshID string) ([]GetApplianceResultItem, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting appliances information", map[string]interface{}{
		"endpointID":    endpointID,
		"serviceMeshID": serviceMeshID,
	})

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
		tflog.Error(ctx, "Failed to encode appliances request body", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return []GetApplianceResultItem{}, fmt.Errorf("failed to encode request bodys: %w", err)
	}

	resp := GetApplianceResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/appliances/query", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create appliances POST request", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return []GetApplianceResultItem{}, fmt.Errorf("failed to POST create requests: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send appliances POST request", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return []GetApplianceResultItem{}, fmt.Errorf("failed to send POST requests: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse appliances HTTP response", map[string]interface{}{
			"error":         err.Error(),
			"endpointID":    endpointID,
			"serviceMeshID": serviceMeshID,
		})
		return []GetApplianceResultItem{}, fmt.Errorf("failed to parse HTTP responses: %w", err)
	}

	tflog.Debug(ctx, "Appliances retrieved successfully", map[string]interface{}{
		"count":         len(resp.Items),
		"endpointID":    endpointID,
		"serviceMeshID": serviceMeshID,
	})
	return resp.Items, nil
}
