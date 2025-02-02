// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppEngineStartStopResult represents the result of an App Engine start or stop operation.
type AppEngineStartStopResult struct {
	Result string `json:"result"`
}

// AppEngineStart sends a request to start the App Engine component and returns the resulting AppEngineStartStopResult
// object. Returns an error if the request fails or the response cannot be parsed.
func AppEngineStart(c *Client) (AppEngineStartStopResult, error) {

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/components/appengine?action=start", c.HostURL), nil)
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

// AppEngineStop sends a request to stop the App Engine component and returns the resulting AppEngineStartStopResult
// object. Returns an error if the request fails or the response cannot be parsed.
func AppEngineStop(c *Client) (AppEngineStartStopResult, error) {

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/components/appengine?action=stop", c.HostURL), nil)
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

// GetAppEngineStatus sends a GET request to retrieve the current status of the App Engine component and returns the
// resulting AppEngineStartStopResult object. Returns an error if the request fails or the response cannot be parsed.
func GetAppEngineStatus(c *Client) (AppEngineStartStopResult, error) {

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/components/appengine/status", c.HostURL), nil)
	if err != nil {
		return resp, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doAdminRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send GET request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return resp, nil
}
