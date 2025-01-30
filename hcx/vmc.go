// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"
)

type SDDC struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	CloudName        string `json:"cloudName,omitempty"`
	CloudURL         string `json:"cloudUrl,omitempty"`
	CloudType        string `json:"cloudType,omitempty"`
	CloudID          string `json:"cloudId,omitempty"`
	ActivationKey    string `json:"activationKey,omitempty"`
	SubscriptionID   string `json:"subscriptionId,omitempty"`
	ActivationStatus string `json:"activationStatus,omitempty"`
	DeploymentStatus string `json:"deploymentStatus,omitempty"`
	State            string `json:"state"`
}

type GetSddcsResults struct {
	SDDCs []SDDC `json:"sddcs"`
}

type VmcAccessToken struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refreshToken"`
}

type CloudAuthorizationBody struct {
	Token string `json:"token"`
}

type ActivateHcxOnSDDCResults struct {
	JobID string `json:"jobId"`
}

type DeactivateHcxOnSDDCResults struct {
	JobID string `json:"jobId"`
}

func VmcAuthenticate(token string) (string, error) {

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.VmcAuthURL,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/api-tokens/authorize?refresh_token=%s", c.HostURL, token), nil)
	if err != nil {
		return "", err
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		return "", err
	}

	resp := VmcAccessToken{}
	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// parse response header

	log.Printf("**************************")
	log.Printf("[Access token] = %+v", resp.AccessToken)
	log.Printf("**************************")

	return resp.AccessToken, nil

}

func CloudAuthenticate(client *Client, token string) error {

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudAuthURL,
	}

	body := CloudAuthorizationBody{
		Token: token,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sessions", c.HostURL), &buf)
	if err != nil {
		return err
	}

	resp, _, err := c.doVmcRequest(req)
	if err != nil {
		return err
	}

	auth := resp.Header.Get("x-hm-authorization")
	if auth == "" {
		return errors.New("cannot authorize hcx cloud")
	}

	// parse response header
	client.HcxToken = auth

	// parse response header

	log.Printf("**************************")
	log.Printf("[Hcx token] = %+v", client.HcxToken)
	log.Printf("**************************")
	return nil

}

func GetSddcByName(client *Client, sddcName string) (SDDC, error) {

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/sddcs", c.HostURL), nil)
	if err != nil {
		return SDDC{}, err
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		return SDDC{}, err
	}

	resp := GetSddcsResults{}
	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return SDDC{}, err
	}

	for _, j := range resp.SDDCs {
		if j.Name == sddcName {
			return j, nil
		}
	}

	// parse response header
	return SDDC{}, errors.New("cant find the sddc")

}

func GetSddcByID(client *Client, sddcID string) (SDDC, error) {

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/sddcs", c.HostURL), nil)
	if err != nil {
		return SDDC{}, err
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		return SDDC{}, err
	}

	resp := GetSddcsResults{}
	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return SDDC{}, err
	}

	for _, j := range resp.SDDCs {
		if j.ID == sddcID {
			return j, nil
		}
	}

	// parse response header
	return SDDC{}, errors.New("cant find the sddc")

}

func ActivateHcxOnSDDC(client *Client, sddcID string) (ActivateHcxOnSDDCResults, error) {

	resp := ActivateHcxOnSDDCResults{}

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sddcs/%s?action=activate", c.HostURL, sddcID), nil)
	if err != nil {
		return resp, err
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		return resp, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// parse response header
	return resp, nil

}

func DeactivateHcxOnSDDC(client *Client, sddcID string) (DeactivateHcxOnSDDCResults, error) {

	resp := DeactivateHcxOnSDDCResults{}

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sddcs/%s?action=deactivate", c.HostURL, sddcID), nil)
	if err != nil {
		return resp, err
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		return resp, err
	}

	// Parse response body.
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// parse response header
	return resp, nil

}
