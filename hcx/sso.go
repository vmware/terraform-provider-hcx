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

type InsertSSOBody struct {
	Data InsertSSOData `json:"data"`
}

type InsertSSOData struct {
	Items []InsertSSODataItem `json:"items"`
}

type InsertSSODataItem struct {
	Config InsertSSODataItemConfig `json:"config"`
}

type InsertSSODataItemConfig struct {
	LookupServiceUrl string `json:"lookupServiceUrl"`
	ProviderType     string `json:"providerType"`
	UUID             string `json:"UUID,omitempty"`
}

type InsertSSOResult struct {
	InsertSSOData InsertSSOData `json:"data"`
}

type DeleteSSOResult struct {
	InsertSSOData InsertSSOData `json:"data"`
}

type GetSSOResult struct {
	InsertSSOData InsertSSOData `json:"data"`
}

// GetSSO ...
func GetSSO(c *Client) (GetSSOResult, error) {

	resp := GetSSOResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice", c.HostURL), nil)
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

// InsertSSO ...
func InsertSSO(c *Client, body InsertSSOBody) (InsertSSOResult, error) {

	resp := InsertSSOResult{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice", c.HostURL), buf)
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

// UpdateSSO ...
func UpdateSSO(c *Client, body InsertSSOBody) (InsertSSOResult, error) {

	resp := InsertSSOResult{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice/%s", c.HostURL, body.Data.Items[0].Config.UUID), buf)
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

// DeleteSSO ...
func DeleteSSO(c *Client, SSOUUID string) (DeleteSSOResult, error) {

	resp := DeleteSSOResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/lookupservice/%s", c.HostURL, SSOUUID), nil)
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
