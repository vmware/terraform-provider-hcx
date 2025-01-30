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

// ActivateBody represents the structure of the request body used for activation actions.
type ActivateBody struct {
	Data ActivateData `json:"data"`
}

// ActivateData represents the detailed activation data, which includes a list of activation items.
type ActivateData struct {
	Items []ActivateDataItem `json:"items"`
}

// ActivateDataItem represents an individual activation item, containing its specific configuration details.
type ActivateDataItem struct {
	Config ActivateDataItemConfig `json:"config"`
}

// ActivateDataItemConfig represents the configuration details for a specific activation item.
type ActivateDataItemConfig struct {
	URL           string `json:"url"`
	ActivationKey string `json:"activationKey"`
	UUID          string `json:"UUID,omitempty"`
}

// PostActivate sends a request to activate a configuration using the provided body and returns the resulting
// ActivateBody object. Returns an error if the request fails or the response cannot be parsed.
func PostActivate(c *Client, body ActivateBody) (ActivateBody, error) {

	resp := ActivateBody{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/hcx", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request.
	_, r, err := c.doAdminRequest(req)
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

// GetActivate sends a request to retrieve the current activation configuration and returns the resulting ActivateBody
// object. Returns an error if the request fails or the response cannot be parsed.
func GetActivate(c *Client) (ActivateBody, error) {

	resp := ActivateBody{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/hcx", c.HostURL), nil)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request.
	_, r, err := c.doAdminRequest(req)
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

// DeleteActivate sends a  request to remove the activation configuration using the provided body and returns the
// resulting ActivateBody object. Returns an error if the request fails or the response cannot be parsed.

func DeleteActivate(c *Client, body ActivateBody) (ActivateBody, error) {

	resp := ActivateBody{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/hcx", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request.
	_, r, err := c.doAdminRequest(req)
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
