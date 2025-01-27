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

type ComputeProfile struct {
	ComputeProfileID   string `json:"computeProfileId"`
	ComputeProfileName string `json:"computeProfileName"`
	EndpointID         string `json:"endpointId"`
	EndpointName       string `json:"endpointName"`
}

type WanoptConfig struct {
	UplinkMaxBandwidth int `json:"uplinkMaxBandwidth"`
}

type TrafficEnggCfg struct {
	IsAppPathResiliencyEnabled   bool `json:"isAppPathResiliencyEnabled"`
	IsTCPFlowConditioningEnabled bool `json:"isTcpFlowConditioningEnabled"`
}

type SwitchPairCount struct {
	Switches          []Switch `json:"switches"`
	L2cApplianceCount int      `json:"l2cApplianceCount"`
}

type InsertServiceMeshBody struct {
	Name            string            `json:"name"`
	ComputeProfiles []ComputeProfile  `json:"computeProfiles"`
	WanoptConfig    WanoptConfig      `json:"wanoptConfig"`
	TrafficEnggCfg  TrafficEnggCfg    `json:"trafficEnggCfg"`
	Services        []Service         `json:"services"`
	SwitchPairCount []SwitchPairCount `json:"switchPairCount"`
}

type InsertServiceMeshResult struct {
	Data InsertServiceMeshData `json:"data"`
}

type InsertServiceMeshData struct {
	InterconnectID string `json:"interconnectTaskId"`
	ServiceMeshID  string `json:"serviceMeshId"`
}

type DeleteServiceMeshResult struct {
	Data DeleteServiceMeshData `json:"data"`
}

type DeleteServiceMeshData struct {
	InterconnectTaskID string `json:"interconnectTaskId"`
	ServiceMeshID      string `json:"serviceMeshId"`
}

// InsertServiceMesh ...
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

// DeleteServiceMesh ...
func DeleteServiceMesh(c *Client, serviceMeshID string, force bool) (DeleteServiceMeshResult, error) {

	resp := DeleteServiceMeshResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/interconnect/serviceMesh/%s?force=%v", c.HostURL, serviceMeshID, force), nil)
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
