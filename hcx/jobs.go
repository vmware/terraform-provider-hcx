// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// AppEngineStartStopResult represents the result of an App Engine start or stop operation.
type AppEngineStartStopResult struct {
	Result string `json:"result"`
}

// AppEngineStart sends a request to start the App Engine component and returns the resulting AppEngineStartStopResult
// object. Returns an error if the request fails or the response cannot be parsed.
func AppEngineStart(c *Client) (AppEngineStartStopResult, error) {
	ctx := context.Background()
	tflog.Info(ctx, "Starting App Engine component")

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/components/appengine?action=start", c.HostURL), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create App Engine start POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send App Engine start POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse App Engine start HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Info(ctx, "App Engine start initiated", map[string]interface{}{
		"result": resp.Result,
	})
	return resp, nil
}

// AppEngineStop sends a request to stop the App Engine component and returns the resulting AppEngineStartStopResult
// object. Returns an error if the request fails or the response cannot be parsed.
func AppEngineStop(c *Client) (AppEngineStartStopResult, error) {
	ctx := context.Background()
	tflog.Info(ctx, "Stopping App Engine component")

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/components/appengine?action=stop", c.HostURL), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create App Engine stop POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send App Engine stop POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse App Engine stop HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Info(ctx, "App Engine stop initiated", map[string]interface{}{
		"result": resp.Result,
	})
	return resp, nil
}

// GetAppEngineStatus sends a GET request to retrieve the current status of the App Engine component and returns the
// resulting AppEngineStartStopResult object. Returns an error if the request fails or the response cannot be parsed.
func GetAppEngineStatus(c *Client) (AppEngineStartStopResult, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Getting App Engine status")

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/components/appengine/status", c.HostURL), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create App Engine status GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send App Engine status GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse App Engine status HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Debug(ctx, "App Engine status retrieved", map[string]interface{}{
		"status": resp.Result,
	})
	return resp, nil
}
