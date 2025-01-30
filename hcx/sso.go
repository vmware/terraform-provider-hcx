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

	resp := GetSSOResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice", c.HostURL), nil)
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

// InsertSSO sends a POST request to create a new SSO configuration using the provided body and returns the resulting
// InsertSSOResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertSSO(c *Client, body InsertSSOBody) (InsertSSOResult, error) {

	resp := InsertSSOResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice", c.HostURL), &buf)
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

// UpdateSSO sends a POST request to update the existing SSO configuration using the provided body. It returns the
// resulting InsertSSOResult object or an error if the request fails or the response cannot be parsed.
func UpdateSSO(c *Client, body InsertSSOBody) (InsertSSOResult, error) {

	resp := InsertSSOResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice/%s", c.HostURL, body.Data.Items[0].Config.UUID), &buf)
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

// DeleteSSO sends a DELETE request to remove the SSO configuration identified by the provided SSOUUID and returns the
// resulting DeleteSSOResult object. Returns an error if the request fails or the response cannot be parsed.
func DeleteSSO(c *Client, SSOUUID string) (DeleteSSOResult, error) {

	resp := DeleteSSOResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice/%s", c.HostURL, SSOUUID), nil)
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
