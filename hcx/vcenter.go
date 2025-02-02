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

// InsertvCenterBody represents the request body for adding a new vCenter instance configuration.
type InsertvCenterBody struct {
	Data InsertvCenterData `json:"data"`
}

// InsertvCenterData represents a collection of vCenter data items required for configuration or operations.
type InsertvCenterData struct {
	Items []InsertvCenterDataItem `json:"items"`
}

// InsertvCenterDataItem represents a single item containing configuration details for a vCenter instance.
type InsertvCenterDataItem struct {
	Config InsertvCenterDataItemConfig `json:"config"`
}

// InsertvCenterDataItemConfig represents the configuration details for connecting to a vCenter instance.
type InsertvCenterDataItemConfig struct {
	URL      string `json:"url"`
	Username string `json:"userName"`
	Password string `json:"password"`
	VcUUID   string `json:"vcuuid,omitempty"`
	UUID     string `json:"UUID,omitempty"`
}

// InsertvCenterResult represents the result returned after inserting a vCenter configuration.
type InsertvCenterResult struct {
	InsertvCenterData InsertvCenterData `json:"data"`
}

// DeletevCenterResult represents the result returned after a vCenter instance is deleted, containing related data.
type DeletevCenterResult struct {
	InsertvCenterData InsertvCenterData `json:"data"`
}

// InsertvCenter sends a request to add a new vCenter instance configuration using the provided body and returns the
// resulting InsertvCenterResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertvCenter(c *Client, body InsertvCenterBody) (InsertvCenterResult, error) {
	resp := InsertvCenterResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/vcenter", c.HostURL), &buf)
	if err != nil {
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}

// DeletevCenter sends a request to remove a vCenter instance configuration identified by the provided vCenterUUID and
// returns the resulting DeletevCenterResult object. Returns an error if the request fails or the response cannot be
// parsed.
func DeletevCenter(c *Client, vCenterUUID string) (DeletevCenterResult, error) {
	resp := DeletevCenterResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/vcenter/%s", c.HostURL, vCenterUUID), nil)
	if err != nil {
		return resp, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send DELETE request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}
