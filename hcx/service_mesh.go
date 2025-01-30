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

// ComputeProfile represents a configuration profile for a specific compute endpoint.
type ComputeProfile struct {
	ComputeProfileID   string `json:"computeProfileId"`
	ComputeProfileName string `json:"computeProfileName"`
	EndpointID         string `json:"endpointId"`
	EndpointName       string `json:"endpointName"`
}

// WanoptConfig represents WAN optimization configuration with uplink maximum bandwidth properties.
type WanoptConfig struct {
	UplinkMaxBandwidth int `json:"uplinkMaxBandwidth"`
}

// TrafficEnggCfg represents the configuration for traffic engineering settings.
type TrafficEnggCfg struct {
	IsAppPathResiliencyEnabled   bool `json:"isAppPathResiliencyEnabled"`
	IsTCPFlowConditioningEnabled bool `json:"isTcpFlowConditioningEnabled"`
}

// SwitchPairCount represents a pairing of switches with a count of Layer 2 appliances in a network configuration.
type SwitchPairCount struct {
	Switches          []Switch `json:"switches"`
	L2cApplianceCount int      `json:"l2cApplianceCount"`
}

// InsertServiceMeshBody represents the body structure required to insert a service mesh configuration.
type InsertServiceMeshBody struct {
	Name            string            `json:"name"`
	ComputeProfiles []ComputeProfile  `json:"computeProfiles"`
	WanoptConfig    WanoptConfig      `json:"wanoptConfig"`
	TrafficEnggCfg  TrafficEnggCfg    `json:"trafficEnggCfg"`
	Services        []Service         `json:"services"`
	SwitchPairCount []SwitchPairCount `json:"switchPairCount"`
}

// InsertServiceMeshResult represents the result returned after inserting a service mesh configuration. It contains the
// data defining the created service mesh.
type InsertServiceMeshResult struct {
	Data InsertServiceMeshData `json:"data"`
}

// InsertServiceMeshData represents the structure for storing information about a service mesh insertion task.
type InsertServiceMeshData struct {
	InterconnectID string `json:"interconnectTaskId"`
	ServiceMeshID  string `json:"serviceMeshId"`
}

// DeleteServiceMeshResult represents the result obtained after attempting to delete a service mesh.
type DeleteServiceMeshResult struct {
	Data DeleteServiceMeshData `json:"data"`
}

// DeleteServiceMeshData represents the payload required to request the deletion of a service mesh.
type DeleteServiceMeshData struct {
	InterconnectTaskID string `json:"interconnectTaskId"`
	ServiceMeshID      string `json:"serviceMeshId"`
}

// InsertServiceMesh sends a request to create a new service mesh using the provided body and returns the resulting
// InsertServiceMeshResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertServiceMesh(c *Client, body InsertServiceMeshBody) (InsertServiceMeshResult, error) {

	resp := InsertServiceMeshResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/serviceMesh", c.HostURL), &buf)
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

// DeleteServiceMesh sends a request to remove a service mesh identified by the serviceMeshID. The force parameter
// determines whether to forcibly delete it. Returns the resulting DeleteServiceMeshResult object or an error if the
// request fails or the response cannot be parsed.
func DeleteServiceMesh(c *Client, serviceMeshID string, force bool) (DeleteServiceMeshResult, error) {

	resp := DeleteServiceMeshResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/interconnect/serviceMesh/%s?force=%v", c.HostURL, serviceMeshID, force), nil)
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
