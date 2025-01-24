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

type SetLocationBody struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	CityAscii string  `json:"cityAscii"`
	Province  string  `json:"province"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GetLocationResult struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Province  string  `json:"province"`
	CityAscii string  `json:"cityAscii"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// SetLocation ...
func SetLocation(c *Client, body SetLocationBody) error {

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s:9443/api/admin/global/config/location", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Send the request
	_, _, err = c.doAdminRequest(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// GetLocation ...
func GetLocation(c *Client) (GetLocationResult, error) {

	resp := GetLocationResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/location", c.HostURL), nil)
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
