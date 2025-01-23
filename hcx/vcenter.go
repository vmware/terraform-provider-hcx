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

type InsertvCenterBody struct {
	Data InsertvCenterData `json:"data"`
}

type InsertvCenterData struct {
	Items []InsertvCenterDataItem `json:"items"`
}

type InsertvCenterDataItem struct {
	Config InsertvCenterDataItemConfig `json:"config"`
}

type InsertvCenterDataItemConfig struct {
	URL      string `json:"url"`
	Username string `json:"userName"`
	Password string `json:"password"`
	Vcuuid   string `json:"vcuuid,omitempty"`
	UUID     string `json:"UUID,omitempty"`
}

type InsertvCenterResult struct {
	InsertvCenterData InsertvCenterData `json:"data"`
}

type DeletevCenterResult struct {
	InsertvCenterData InsertvCenterData `json:"data"`
}

// InsertvCenter ...
func InsertvCenter(c *Client, body InsertvCenterBody) (InsertvCenterResult, error) {

	resp := InsertvCenterResult{}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s:9443/api/admin/global/config/vcenter", c.HostURL), &buf)
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

// DeletevCenter ...
func DeletevCenter(c *Client, vCenterUUID string) (DeletevCenterResult, error) {

	resp := DeletevCenterResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s:9443/api/admin/global/config/vcenter/%s", c.HostURL, vCenterUUID), nil)
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
