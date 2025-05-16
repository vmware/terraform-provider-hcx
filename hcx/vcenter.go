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
	ctx := context.Background()
	tflog.Debug(ctx, "Creating new vCenter instance", map[string]interface{}{
		"url":      body.Data.Items[0].Config.URL,
		"username": body.Data.Items[0].Config.Username,
	})

	resp := InsertvCenterResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode vCenter request body", map[string]interface{}{
			"error": err.Error(),
			"url":   body.Data.Items[0].Config.URL,
		})
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/vcenter", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create vCenter POST request", map[string]interface{}{
			"error": err.Error(),
			"url":   body.Data.Items[0].Config.URL,
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send vCenter POST request", map[string]interface{}{
			"error": err.Error(),
			"url":   body.Data.Items[0].Config.URL,
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse vCenter HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "vCenter instance created successfully", map[string]interface{}{
		"uuid": resp.InsertvCenterData.Items[0].Config.UUID,
		"url":  body.Data.Items[0].Config.URL,
	})
	return resp, nil
}

// DeletevCenter sends a request to remove a vCenter instance configuration identified by the provided vCenterUUID and
// returns the resulting DeletevCenterResult object. Returns an error if the request fails or the response cannot be
// parsed.
func DeletevCenter(c *Client, vCenterUUID string) (DeletevCenterResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Deleting vCenter instance", map[string]interface{}{
		"uuid": vCenterUUID,
	})

	resp := DeletevCenterResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/vcenter/%s", c.HostURL, vCenterUUID), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create vCenter DELETE request", map[string]interface{}{
			"error": err.Error(),
			"uuid":  vCenterUUID,
		})
		return resp, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send vCenter DELETE request", map[string]interface{}{
			"error": err.Error(),
			"uuid":  vCenterUUID,
		})
		return resp, fmt.Errorf("failed to send DELETE request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse vCenter DELETE HTTP response", map[string]interface{}{
			"error": err.Error(),
			"uuid":  vCenterUUID,
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "vCenter instance deleted successfully", map[string]interface{}{
		"uuid": vCenterUUID,
	})
	return resp, nil
}
