// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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

type Compute struct {
	CmpId   string `json:"cmpId"`
	CmpName string `json:"cmpName"`
	CmpType string `json:"cmpType"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

type Storage struct {
	CmpId   string `json:"cmpId"`
	CmpName string `json:"cmpName"`
	CmpType string `json:"cmpType"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

type DeploymentContainer struct {
	Computes          []Compute `json:"compute"`
	CpuReservation    int       `json:"cpuReservation"`
	MemoryReservation int       `json:"memoryReservation"`
	Storage           []Storage `json:"storage"`
}

type Network struct {
	Name         string        `json:"name"`
	ID           string        `json:"id"`
	StaticRoutes []interface{} `json:"staticRoutes"`
	Status       Status        `json:"status"`
	Tags         []string      `json:"tags"`
}

type Status struct {
	State string `json:"state"`
}

type Service struct {
	Name string `json:"name"`
}

type Switch struct {
	CmpID  string `json:"cmpId"`
	ID     string `json:"id"`
	MaxMTU int    `json:"maxMtu,omitempty"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type InsertComputeProfileResult struct {
	Data InsertComputeProfileResultData `json:"data"`
}

type InsertComputeProfileResultData struct {
	InterconnectTaskId string `json:"interconnectTaskId"`
	ComputeProfileId   string `json:"computeProfileId"`
}

type GetComputeProfileResult struct {
	Items []GetComputeProfileResultItem `json:"items"`
}

type GetComputeProfileResultItem struct {
	ComputeProfileId     string              `json:"computeProfileId"`
	Name                 string              `json:"name"`
	Compute              []Compute           `json:"compute"`
	Services             []Service           `json:"services"`
	DeploymentContainers DeploymentContainer `json:"deploymentContainer"`
	Networks             []Network           `json:"networks"`
	State                string              `json:"state"`
	Switches             []Switch            `json:"switches"`
}

// InsertComputeProfile ...
func InsertComputeProfile(c *Client, body InsertComputeProfileBody) (InsertComputeProfileResult, error) {

	resp := InsertComputeProfileResult{}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/interconnect/computeProfiles", c.HostURL), &buf)
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

// DeleteComputeProfile ...
func DeleteComputeProfile(c *Client, computeprofileID string) (InsertComputeProfileResult, error) {

	resp := InsertComputeProfileResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/interconnect/computeProfiles/%s", c.HostURL, computeprofileID), nil)
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

// GetComputeProfile ...
func GetComputeProfile(c *Client, endpointId string, computeprofileName string) (GetComputeProfileResultItem, error) {

	resp := GetComputeProfileResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/interconnect/computeProfiles?endpointId=%s", c.HostURL, endpointId), nil)
	if err != nil {
		fmt.Println(err)
		return GetComputeProfileResultItem{}, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetComputeProfileResultItem{}, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetComputeProfileResultItem{}, err
	}

	for _, j := range resp.Items {
		if j.Name == computeprofileName {
			return j, nil
		}
	}

	return GetComputeProfileResultItem{}, errors.New("cant find compute profile")
}
