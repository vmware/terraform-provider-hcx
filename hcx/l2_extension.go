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

// InsertL2ExtensionBody represents the request body structure for creating a Layer 2 extension.
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

// DestinationNetwork represents the configuration for a destination network in a Layer 2 extension setup.
type DestinationNetwork struct {
	GatewayID string `json:"gatewayId"`
}

// Destination represents the primary structure for encapsulating endpoint and resource details in a network extension
// setup.
type Destination struct {
	EndpointID   string `json:"endpointId"`
	EndpointName string `json:"endpointName"`
	EndpointType string `json:"endpointType"`
	ResourceID   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
}

// Features defines a struct for enabling specific configurations such as egress optimization and mobility-optimized
// networking.
type Features struct {
	EgressOptimization bool `json:"egressOptimization"`
	Mon                bool `json:"mobilityOptimizedNetworking"`
}

// SourceAppliance represents the source appliance information in a Layer 2 network extension configuration.
type SourceAppliance struct {
	ApplianceID string `json:"applianceId"`
}

// SourceNetwork represents the details of a network within a Layer 2 extension configuration.
type SourceNetwork struct {
	NetworkID   string `json:"networkId"`
	NetworkName string `json:"networkName"`
	NetworkType string `json:"networkType"`
}

// InsertL2ExtensionResult represents the result of an InsertL2Extension operation.
type InsertL2ExtensionResult struct {
	ID string `json:"id"`
}

// GetL2ExtensionsResult represents the result of a request for fetching Layer 2 extensions information.
type GetL2ExtensionsResult struct {
	Items []GetL2ExtensionsResultItem `json:"items"`
}

// GetL2ExtensionsResultItem represents an item in the result of a Layer 2 extensions request.
type GetL2ExtensionsResultItem struct {
	StretchID       string          `json:"stretchId"`
	OperationStatus OperationStatus `json:"operationStatus"`
	SourceNetwork   SourceNetwork   `json:"sourceNetwork"`
}

// OperationStatus represents the status of an operation.
type OperationStatus struct {
	State string `json:"state"`
}

// DeleteL2ExtensionResult represents the result of a successful deletion of an L2 extension.
type DeleteL2ExtensionResult struct {
	ID string `json:"id"`
}

// InsertL2Extension sends a POST request to create a new L2 extension using the provided body and returns the resulting
// InsertL2ExtensionResult object. Returns an error if the request fails or the response cannot be parsed.
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

// GetL2Extensions sends a GET request to retrieve a list of L2 extensions and returns the resulting
// GetL2ExtensionsResultItem object matching the given networkName. Returns an error if the request fails, the response
// cannot be parsed, or no matching L2 extension is found.
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

// DeleteL2Extension sends a DELETE request to remove an L2 extension with the provided stretchID and returns the
// resulting DeleteL2ExtensionResult object. Returns an error if the request fails or the response cannot be parsed.
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
