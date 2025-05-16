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

// RemoteData represents the structure required for remote connection configurations.
type RemoteData struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	URL        string `json:"url"`
	EndpointID string `json:"endpointId,omitempty"`
	CloudType  string `json:"cloudType,omitempty"`
}

// RemoteCloudConfigBody represents the request body for configuring remote cloud connections.
type RemoteCloudConfigBody struct {
	Remote RemoteData `json:"remote"`
}

// PostRemoteCloudConfigResultData represents the result data containing the job identifier for a remote cloud
// configuration task.
type PostRemoteCloudConfigResultData struct {
	JobID string `json:"jobId"`
}

// PostRemoteCloudConfigResult represents the result of posting a remote cloud configuration, including metadata and
// outcomes.
type PostRemoteCloudConfigResult struct {
	Success   bool                               `json:"success"`
	Completed bool                               `json:"completed"`
	Time      int64                              `json:"time"`
	Version   string                             `json:"version"`
	Data      PostRemoteCloudConfigResultData    `json:"data"`
	Errors    []PostRemoteCloudConfigResultError `json:"errors"`
}

// PostRemoteCloudConfigResultError represents an error from posting a remote cloud configuration.
type PostRemoteCloudConfigResultError struct {
	Error string                   `json:"error"`
	Text  string                   `json:"text"`
	Data  []map[string]interface{} `json:"data"`
}

// GetRemoteCloudConfigResult represents the result of retrieving remote cloud configuration data.
type GetRemoteCloudConfigResult struct {
	Success   bool                           `json:"success"`
	Completed bool                           `json:"completed"`
	Time      int64                          `json:"time"`
	Version   string                         `json:"version"`
	Data      GetRemoteCloudConfigResultData `json:"data"`
}

// GetRemoteCloudConfigResultData represents the result data containing remote cloud configuration items.
type GetRemoteCloudConfigResultData struct {
	Items []RemoteData `json:"items"`
}

// DeleteRemoteCloudConfigResult represents the outcome of a request to delete a remote cloud configuration.
// It includes success status, completion status, and the time of the operation as a Unix timestamp.
type DeleteRemoteCloudConfigResult struct {
	Success   bool  `json:"success"`
	Completed bool  `json:"completed"`
	Time      int64 `json:"time"`
}

// InsertSitePairing sends a request to create a new site pairing using the provided body and returns the resulting
// PostRemoteCloudConfigResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertSitePairing(c *Client, body RemoteCloudConfigBody) (PostRemoteCloudConfigResult, error) {
	resp := PostRemoteCloudConfigResult{}
	ctx := context.Background()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		tflog.Error(ctx, "Failed to encode site pairing request body", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	tflog.Debug(ctx, "Creating site pairing request", map[string]interface{}{
		"url": fmt.Sprintf("%s/hybridity/api/cloudConfigs", c.HostURL),
	})

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/cloudConfigs", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create site pairing request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	req = req.WithContext(ctx)

	_, r, err := c.doRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send site pairing request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse site pairing response", map[string]interface{}{
			"error": err.Error(),
			"body":  string(r),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "Site pairing request successful", map[string]interface{}{
		"success":    resp.Success,
		"completed":  resp.Completed,
		"has_errors": resp.Errors != nil && len(resp.Errors) > 0,
	})

	return resp, nil
}

// GetSitePairings sends a GET request to retrieve all existing site pairings and returns the resulting
// GetRemoteCloudConfigResult object. Returns an error if the request fails or the response cannot be parsed.
func GetSitePairings(c *Client) (GetRemoteCloudConfigResult, error) {
	resp := GetRemoteCloudConfigResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/cloudConfigs", c.HostURL), nil)
	if err != nil {
		return resp, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}

// DeleteSitePairings sends a DELETE request to remove a site pairing identified by the provided endpointID and returns
// the resulting DeleteRemoteCloudConfigResult object. Returns an error if the request fails or the response cannot be
// parsed.
func DeleteSitePairings(c *Client, endpointID string) (DeleteRemoteCloudConfigResult, error) {
	resp := DeleteRemoteCloudConfigResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/endpointPairing/%s", c.HostURL, endpointID), nil)
	if err != nil {
		return resp, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send DELETE request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}
