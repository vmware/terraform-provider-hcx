// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// InsertComputeProfileBody represents the body structure for inserting a compute profile.
type InsertComputeProfileBody struct {
	Computes             []Compute           `json:"compute"`
	ComputeProfileID     string              `json:"computeProfileId"`
	DeploymentContainers DeploymentContainer `json:"deploymentContainer"`
	Name                 string              `json:"name"`
	Networks             []Network           `json:"networks"`
	Services             []Service           `json:"services"`
	State                string              `json:"state"`
	Switches             []Switch            `json:"switches"`
}

// Compute represents the compute resource configuration.
type Compute struct {
	ComputeID   string `json:"cmpId"`
	ComputeName string `json:"cmpName"`
	ComputeType string `json:"cmpType"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

// Storage represents a storage entity in a compute profile configuration.
type Storage struct {
	ComputeID   string `json:"cmpId"`
	ComputeName string `json:"cmpName"`
	ComputeType string `json:"cmpType"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

// DeploymentContainer represents a container holding deployment configuration.
type DeploymentContainer struct {
	Computes          []Compute `json:"compute"`
	CPUReservation    int       `json:"cpuReservation"`
	MemoryReservation int       `json:"memoryReservation"`
	Storage           []Storage `json:"storage"`
}

// Network represents a network entity with associated details.
type Network struct {
	Name         string        `json:"name"`
	ID           string        `json:"id"`
	StaticRoutes []interface{} `json:"staticRoutes"`
	Status       Status        `json:"status"`
	Tags         []string      `json:"tags"`
}

// Status represents the current state of an entity, typically used to indicate its operational or lifecycle state.
type Status struct {
	State string `json:"state"`
}

// Service represents a service with a specific name in a system configuration.
type Service struct {
	Name string `json:"name"`
}

// Switch represents a network switch with identifiable attributes.
type Switch struct {
	ComputeID string `json:"cmpId"`
	ID        string `json:"id"`
	MaxMTU    int    `json:"maxMtu,omitempty"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

// InsertComputeProfileResult represents the result of an operation to insert a compute profile.
type InsertComputeProfileResult struct {
	Data InsertComputeProfileResultData `json:"data"`
}

// InsertComputeProfileResultData represents the result data for inserting a compute profile.
type InsertComputeProfileResultData struct {
	InterconnectTaskID string `json:"interconnectTaskId"`
	ComputeProfileID   string `json:"computeProfileId"`
}

// GetComputeProfileResult represents a collection of compute profile details.
type GetComputeProfileResult struct {
	Items []GetComputeProfileResultItem `json:"items"`
}

// GetComputeProfileResultItem represents the details of a compute profile.
type GetComputeProfileResultItem struct {
	ComputeProfileID     string              `json:"computeProfileId"`
	Name                 string              `json:"name"`
	Compute              []Compute           `json:"compute"`
	Services             []Service           `json:"services"`
	DeploymentContainers DeploymentContainer `json:"deploymentContainer"`
	Networks             []Network           `json:"networks"`
	State                string              `json:"state"`
	Switches             []Switch            `json:"switches"`
}

// InsertComputeProfile sends a request to create a new compute profile using the provided body and returns an
// InsertComputeProfileResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertComputeProfile(c *Client, body InsertComputeProfileBody) (InsertComputeProfileResult, error) {

	resp := InsertComputeProfileResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/computeProfiles", c.HostURL), &buf)
	if err != nil {
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal POST response: %w", err)
	}

	return resp, nil
}

// DeleteComputeProfile sends a request to delete a specific compute profile identified by computeProfileID and an
// InsertComputeProfileResult object indicating the result of the operation. Returns an error if the request fails or
// the response cannot be parsed.
func DeleteComputeProfile(c *Client, computeProfileID string) (InsertComputeProfileResult, error) {

	resp := InsertComputeProfileResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/interconnect/computeProfiles/%s", c.HostURL, computeProfileID), nil)
	if err != nil {
		return resp, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send DELETE request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal DELETE response: %w", err)
	}

	return resp, nil
}

// GetComputeProfile retrieves the details of a compute profile using the provided endpointID and computeProfileName,
// returning a GetComputeProfileResultItem object for the matching profile. Returns an error if the request fails, the
// response cannot be parsed, or no matching profile is found.
func GetComputeProfile(c *Client, endpointID string, computeProfileName string) (GetComputeProfileResultItem, error) {

	resp := GetComputeProfileResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/interconnect/computeProfiles?endpointId=%s", c.HostURL, endpointID), nil)
	if err != nil {
		return GetComputeProfileResultItem{}, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return GetComputeProfileResultItem{}, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return GetComputeProfileResultItem{}, fmt.Errorf("failed to unmarshal GET response: %w", err)
	}

	for _, j := range resp.Items {
		if j.Name == computeProfileName {
			return j, nil
		}
	}

	return GetComputeProfileResultItem{}, fmt.Errorf("cannot find compute profile: %s", computeProfileName)
}
