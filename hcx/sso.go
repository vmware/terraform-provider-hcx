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

// InsertSSOBody represents the structure required for inserting  Single Sign-On configurations with associated data.
type InsertSSOBody struct {
	Data InsertSSOData `json:"data"`
}

// InsertSSOData represents the structure for inserting Single Sign-On configuration data, containing a list of
// configuration items.
type InsertSSOData struct {
	Items []InsertSSODataItem `json:"items"`
}

// InsertSSODataItem represents a single  Single Sign-On data configuration item containing its associated configuration
// details.
type InsertSSODataItem struct {
	Config InsertSSODataItemConfig `json:"config"`
}

// InsertSSODataItemConfig represents the configuration for a Single Sign-On data item.
type InsertSSODataItemConfig struct {
	LookupServiceURL string `json:"lookupServiceUrl"`
	ProviderType     string `json:"providerType"`
	UUID             string `json:"UUID,omitempty"`
}

// InsertSSOResult represents the result of an operation to insert a Single Sign-On configuration data.
type InsertSSOResult struct {
	InsertSSOData InsertSSOData `json:"data"`
}

// DeleteSSOResult represents the result of a delete operation for Single Sign-On.
type DeleteSSOResult struct {
	InsertSSOData InsertSSOData `json:"data"`
}

// GetSSOResult represents the result of retrieving the current Single Sign-On configuration.
type GetSSOResult struct {
	InsertSSOData InsertSSOData `json:"data"`
}

// GetSSO sends a GET request to retrieve the current SSO configuration and returns the resulting GetSSOResult object.
// Returns an error if the request fails or the response cannot be parsed.
func GetSSO(c *Client) (GetSSOResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting SSO configuration")

	resp := GetSSOResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice", c.HostURL), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create SSO GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send SSO GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse SSO HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "Successfully retrieved SSO configuration", map[string]interface{}{
		"itemCount": len(resp.InsertSSOData.Items),
	})
	return resp, nil
}

// InsertSSO sends a POST request to create a new SSO configuration using the provided body and returns the resulting
// InsertSSOResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertSSO(c *Client, body InsertSSOBody) (InsertSSOResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Creating new SSO configuration", map[string]interface{}{
		"url": body.Data.Items[0].Config.LookupServiceURL,
	})

	resp := InsertSSOResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode SSO request body", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create SSO POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send SSO POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse SSO HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "SSO configuration created successfully", map[string]interface{}{
		"uuid": resp.InsertSSOData.Items[0].Config.UUID,
	})
	return resp, nil
}

// UpdateSSO sends a POST request to update the existing SSO configuration using the provided body. It returns the
// resulting InsertSSOResult object or an error if the request fails or the response cannot be parsed.
func UpdateSSO(c *Client, body InsertSSOBody) (InsertSSOResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Updating SSO configuration", map[string]interface{}{
		"uuid": body.Data.Items[0].Config.UUID,
		"url":  body.Data.Items[0].Config.LookupServiceURL,
	})

	resp := InsertSSOResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode SSO update request body", map[string]interface{}{
			"error": err.Error(),
			"uuid":  body.Data.Items[0].Config.UUID,
		})
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice/%s", c.HostURL, body.Data.Items[0].Config.UUID), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create SSO update POST request", map[string]interface{}{
			"error": err.Error(),
			"uuid":  body.Data.Items[0].Config.UUID,
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send SSO update POST request", map[string]interface{}{
			"error": err.Error(),
			"uuid":  body.Data.Items[0].Config.UUID,
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse SSO update HTTP response", map[string]interface{}{
			"error": err.Error(),
			"uuid":  body.Data.Items[0].Config.UUID,
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "SSO configuration updated successfully", map[string]interface{}{
		"uuid": body.Data.Items[0].Config.UUID,
	})
	return resp, nil
}

// DeleteSSO sends a DELETE request to remove the SSO configuration identified by the provided SSOUUID and returns the
// resulting DeleteSSOResult object. Returns an error if the request fails or the response cannot be parsed.
func DeleteSSO(c *Client, SSOUUID string) (DeleteSSOResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Deleting SSO configuration", map[string]interface{}{
		"uuid": SSOUUID,
	})

	resp := DeleteSSOResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice/%s", c.HostURL, SSOUUID), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create SSO DELETE request", map[string]interface{}{
			"error": err.Error(),
			"uuid":  SSOUUID,
		})
		return resp, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send SSO DELETE request", map[string]interface{}{
			"error": err.Error(),
			"uuid":  SSOUUID,
		})
		return resp, fmt.Errorf("failed to send DELETE request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse SSO DELETE HTTP response", map[string]interface{}{
			"error": err.Error(),
			"uuid":  SSOUUID,
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "SSO configuration deleted successfully", map[string]interface{}{
		"uuid": SSOUUID,
	})
	return resp, nil
}
