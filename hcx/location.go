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

// SetLocationBody represents a structured request body for configuring location data in a system.
type SetLocationBody struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	CityASCII string  `json:"cityAscii"`
	Province  string  `json:"province"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetLocationResult represents the response structure for location details fetched from an external service.
type GetLocationResult struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Province  string  `json:"province"`
	CityASCII string  `json:"cityAscii"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// SetLocation sends request to update the location configuration using the provided body. Returns an error if the
// request fails or cannot be sent.
func SetLocation(c *Client, body SetLocationBody) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s:9443/api/admin/global/config/location", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Send the request.
	_, _, err = c.doAdminRequest(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// GetLocation sends a request to retrieve the current location configuration and returns the resulting
// GetLocationResult object. Returns an error if the request fails or the response cannot be parsed.
func GetLocation(c *Client) (GetLocationResult, error) {

	resp := GetLocationResult{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s:9443/api/admin/global/config/location", c.HostURL), nil)
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
