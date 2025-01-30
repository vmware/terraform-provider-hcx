// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppEngineStartStopResult struct {
	Result string `json:"result"`
}

// AppEngineStart ...
func AppEngineStart(c *Client) (AppEngineStartStopResult, error) {

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/components/appengine?action=start", c.HostURL), nil)
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

// AppEngineStop ...
func AppEngineStop(c *Client) (AppEngineStartStopResult, error) {

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/components/appengine?action=stop", c.HostURL), nil)
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

// GetAppEngineStatus ...
func GetAppEngineStatus(c *Client) (AppEngineStartStopResult, error) {

	resp := AppEngineStartStopResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/components/appengine/status", c.HostURL), nil)
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
