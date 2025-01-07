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

type ActivateBody struct {
	Data ActivateData `json:"data"`
}

type ActivateData struct {
	Items []ActivateDataItem `json:"items"`
}

type ActivateDataItem struct {
	Config ActivateDataItemConfig `json:"config"`
}

type ActivateDataItemConfig struct {
	URL           string `json:"url"`
	ActivationKey string `json:"activationKey"`
	UUID          string `json:"UUID,omitempty"`
}

// PostActivate ...
func PostActivate(c *Client, body ActivateBody) (ActivateBody, error) {

	resp := ActivateBody{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/hcx", c.HostURL), buf)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request
	_, r, err := c.doAdminRequest(req)
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

// GetActivate ...
func GetActivate(c *Client) (ActivateBody, error) {

	resp := ActivateBody{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/hcx", c.HostURL), nil)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request
	_, r, err := c.doAdminRequest(req)
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

// DeleteActivate ...
func DeleteActivate(c *Client, body ActivateBody) (ActivateBody, error) {

	resp := ActivateBody{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/hcx", c.HostURL), buf)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request
	_, r, err := c.doAdminRequest(req)
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
