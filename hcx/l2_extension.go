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

type InsertL2ExtensionBody struct {
	VcGUID             string             `json:"vcGuid"`
	Gateway            string             `json:"gateway"`
	Netmask            string             `json:"netmask"`
	DNS                []string           `json:"dns"`
	Destination        Destination        `json:"destination"`
	DestinationNetwork DestinationNetwork `json:"destinationNetwork"`
	Features           Features           `json:"features"`
	SourceAppliance    SourceAppliance    `json:"sourceAppliance"`
	SourceNetwork      SourceNetwork      `json:"sourceNetwork"`
}

type DestinationNetwork struct {
	GatewayID string `json:"gatewayId"`
}

type Destination struct {
	EndpointID   string `json:"endpointId"`
	EndpointName string `json:"endpointName"`
	EndpointType string `json:"endpointType"`
	ResourceID   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
}

type Features struct {
	EgressOptimization bool `json:"egressOptimization"`
	Mon                bool `json:"mobilityOptimizedNetworking"`
}

type SourceAppliance struct {
	ApplianceID string `json:"applianceId"`
}

type SourceNetwork struct {
	NetworkID   string `json:"networkId"`
	NetworkName string `json:"networkName"`
	NetworkType string `json:"networkType"`
}

type InsertL2ExtensionResult struct {
	ID string `json:"id"`
}

type GetL2ExtensionsResult struct {
	Items []GetL2ExtensionsResultItem `json:"items"`
}

type GetL2ExtensionsResultItem struct {
	StretchID       string          `json:"stretchId"`
	OperationStatus OperationStatus `json:"operationStatus"`
	SourceNetwork   SourceNetwork   `json:"sourceNetwork"`
}

type OperationStatus struct {
	State string `json:"state"`
}

type DeleteL2ExtensionResult struct {
	ID string `json:"id"`
}

// InsertL2Extention ...
func InsertL2Extension(c *Client, body InsertL2ExtensionBody) (InsertL2ExtensionResult, error) {

	resp := InsertL2ExtensionResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/l2Extensions", c.HostURL), &buf)
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

// GetL2Extensions ...
func GetL2Extensions(c *Client, networkName string) (GetL2ExtensionsResultItem, error) {

	resp := GetL2ExtensionsResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/l2Extensions", c.HostURL), nil)
	if err != nil {
		fmt.Println(err)
		return GetL2ExtensionsResultItem{}, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return GetL2ExtensionsResultItem{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return GetL2ExtensionsResultItem{}, err
	}

	for _, j := range resp.Items {
		if j.SourceNetwork.NetworkName == networkName {
			return j, nil
		}
	}

	return GetL2ExtensionsResultItem{}, errors.New("cant find compute L2 extension")
}

// DeleteL2Extension ...
func DeleteL2Extension(c *Client, stretchID string) (DeleteL2ExtensionResult, error) {

	resp := DeleteL2ExtensionResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/l2Extensions/%s", c.HostURL, stretchID), nil)
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
