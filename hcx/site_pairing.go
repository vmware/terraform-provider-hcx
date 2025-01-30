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

type RemoteData struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	URL        string `json:"url"`
	EndpointID string `json:"endpointId,omitempty"`
	CloudType  string `json:"cloudType,omitempty"`
}

type RemoteCloudConfigBody struct {
	Remote RemoteData `json:"remote"`
}

type PostRemoteCloudConfigResultData struct {
	JobID string `json:"jobId"`
}

type PostRemoteCloudConfigResult struct {
	Success   bool                               `json:"success"`
	Completed bool                               `json:"completed"`
	Time      int64                              `json:"time"`
	Version   string                             `json:"version"`
	Data      PostRemoteCloudConfigResultData    `json:"data"`
	Errors    []PostRemoteCloudConfigResultError `json:"errors"`
}

type PostRemoteCloudConfigResultError struct {
	Error string                   `json:"error"`
	Text  string                   `json:"text"`
	Data  []map[string]interface{} `json:"data"`
}

type GetRemoteCloudConfigResult struct {
	Success   bool                           `json:"success"`
	Completed bool                           `json:"completed"`
	Time      int64                          `json:"time"`
	Version   string                         `json:"version"`
	Data      GetRemoteCloudConfigResultData `json:"data"`
}

type GetRemoteCloudConfigResultData struct {
	Items []RemoteData `json:"items"`
}

type DeleteRemoteCloudConfigResult struct {
	Success   bool  `json:"success"`
	Completed bool  `json:"completed"`
	Time      int64 `json:"time"`
}

// InsertSitePairing ...
func InsertSitePairing(c *Client, body RemoteCloudConfigBody) (PostRemoteCloudConfigResult, error) {

	resp := PostRemoteCloudConfigResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/cloudConfigs", c.HostURL), &buf)
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

// GetSitePairings ...
func GetSitePairings(c *Client) (GetRemoteCloudConfigResult, error) {

	resp := GetRemoteCloudConfigResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hybridity/api/cloudConfigs", c.HostURL), nil)
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

// DeleteSitePairings ...
func DeleteSitePairings(c *Client, endpointID string) (DeleteRemoteCloudConfigResult, error) {

	resp := DeleteRemoteCloudConfigResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/endpointPairing/%s", c.HostURL, endpointID), nil)
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
