// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type NetworkProfileBody struct {
	Backings        []Backing `json:"backings"`
	Description     string    `json:"description"`
	Organization    string    `json:"organization,omitempty"`
	IPScopes        []IPScope `json:"ipScopes"`
	MTU             int       `json:"mtu"`
	Name            string    `json:"name"`
	L3TenantManaged bool      `json:"l3TenantManaged"`
	OwnedBySystem   bool      `json:"ownedBySystem"`
	ObjectID        string    `json:"objectId,omitempty"`
}

type Filter struct {
	OwnedBySystem        bool `json:"ownedBySystem"`
	AllowTrunkInterfaces bool `json:"allowTrunkInterfaces"`
}

type NetworkFilter struct {
	Filter Filter `json:"filter"`
}

type Backing struct {
	BackingID           string `json:"backingId"`
	BackingName         string `json:"backingName"`
	Type                string `json:"type"`
	VCenterInstanceUUID string `json:"vCenterInstanceUuid"`
	VCenterName         string `json:"vCenterName,omitempty"`
}

type IPScope struct {
	DNSSuffix       string           `json:"dnsSuffix,omitempty"`
	Gateway         string           `json:"gateway,omitempty"`
	PrefixLength    int              `json:"prefixLength"`
	PrimaryDNS      string           `json:"primaryDns,omitempty"`
	SecondaryDNS    string           `json:"secondaryDns,omitempty"`
	NetworkIPRanges []NetworkIPRange `json:"networkIpRanges,omitempty"`
	PoolID          string           `json:"poolId"`
}

type NetworkIPRange struct {
	EndAddress   string `json:"endAddress"`
	StartAddress string `json:"startAddress"`
}

type NetworkProfileResult struct {
	Success   bool               `json:"success"`
	Completed bool               `json:"completed"`
	Time      int64              `json:"time"`
	Data      NetworkProfileData `json:"data"`
}

type NetworkProfileData struct {
	JobID    string `json:"jobId"`
	ObjectID string `json:"objectId"`
}

// InsertNetworkProfile ...
func InsertNetworkProfile(c *Client, body NetworkProfileBody) (NetworkProfileResult, error) {

	resp := NetworkProfileResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/admin/hybridity/api/networks", c.HostURL), &buf)
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

// GetNetworkProfile ...
func GetNetworkProfile(c *Client, name string) (NetworkProfileBody, error) {

	resp := []NetworkProfileBody{}
	body := NetworkFilter{
		Filter: Filter{
			OwnedBySystem:        true,
			AllowTrunkInterfaces: false,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/networks?action=queryIpUsage", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	for _, j := range resp {
		if j.Name == name {
			return j, nil
		}
	}

	return NetworkProfileBody{}, errors.New("cannot find network profile")
}

// GetNetworkProfileByID ...
func GetNetworkProfileByID(c *Client, id string) (NetworkProfileBody, error) {

	resp := []NetworkProfileBody{}
	body := NetworkFilter{
		Filter: Filter{
			OwnedBySystem:        true,
			AllowTrunkInterfaces: false,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/networks?action=queryIpUsage", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	// Send the request.
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return NetworkProfileBody{}, err
	}

	for _, j := range resp {
		if j.ObjectID == id {
			return j, nil
		}
	}

	return NetworkProfileBody{}, errors.New("cannot find network profile")
}

// DeleteNetworkProfile ...
func DeleteNetworkProfile(c *Client, networkID string) (NetworkProfileResult, error) {

	resp := NetworkProfileResult{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hybridity/api/networks/%s", c.HostURL, networkID), nil)
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

// UpdateNetworkProfile ...
func UpdateNetworkProfile(c *Client, body NetworkProfileBody) (NetworkProfileResult, error) {

	resp := NetworkProfileResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/hybridity/api/networks/%s", c.HostURL, body.ObjectID), &buf)
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
