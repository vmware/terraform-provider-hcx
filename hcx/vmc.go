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

// SDDC defines the properties of a Software-Defined Data Center.
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

// GetSddcsResults represents the structure containing a list of Software-Defined Data Centers.
type GetSddcsResults struct {
	SDDCs []SDDC `json:"sddcs"`
}

// VmcAccessToken represents a structure for storing VMware Cloud API authentication tokens and token metadata.
type VmcAccessToken struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refreshToken"`
}

// CloudAuthorizationBody represents the JSON payload containing the authorization token for cloud authentication.
type CloudAuthorizationBody struct {
	Token string `json:"token"`
}

// ActivateHcxOnSDDCResults represents the result of activating HCX on a Software-Defined Data Center.
type ActivateHcxOnSDDCResults struct {
	JobID string `json:"jobId"`
}

// DeactivateHcxOnSDDCResults represents the result of deactivating HCX on a Software-Defined Data Center.
type DeactivateHcxOnSDDCResults struct {
	JobID string `json:"jobId"`
}

// VmcAuthenticate sends a request to authenticate with the VMware Cloud (VMC) API using the provided token. It returns
// an access token as a string or an error if the request fails or the response cannot be parsed.
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

	// Parse response header.

	log.Printf("**************************")
	log.Printf("[Access token] = %+v", resp.AccessToken)
	log.Printf("**************************")

	return resp.AccessToken, nil

}

// CloudAuthenticate sends a request to authenticate to the HCX cloud service using the provided token. On successful
// authentication, it sets the HcxToken field of the provided Client. Returns an error if the request fails, the
// response is invalid, or the token cannot be retrieved.
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

	// Parse response header.
	client.HcxToken = auth

	return nil

}

// GetSddcByName sends a request to retrieve a list of SDDCs (Software-Defined Data Centers) and searches for an SDDC
// with the specified sddcName. It returns the matching SDDC object or an error if the request fails, the response
// cannot be parsed, or the SDDC is not found.
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

	// Parse response header.
	return SDDC{}, errors.New("cant find the sddc")

}

// GetSddcByID sends a request to retrieve a list of SDDCs (Software-Defined Data Centers) and searches for an SDDC with
// the specified sddcID. It returns the matching SDDC object or an error if the request fails, the response cannot be
// parsed, or the SDDC is not found.
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

	// Parse response header.
	return SDDC{}, errors.New("cant find the sddc")

}

// ActivateHcxOnSDDC sends a request to activate HCX on the specified SDDC (Software-Defined Data Center) identified by
// the provided sddcID. It returns the resulting ActivateHcxOnSDDCResults object or an error if the request fails or the
// response cannot be parsed.
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

	// Parse response header.
	return resp, nil

}

// DeactivateHcxOnSDDC sends a POST request to deactivate HCX on the specified SDDC (Software-Defined Data Center)
// identified by the provided sddcID. It returns the resulting DeactivateHcxOnSDDCResults object or an error if the
// request fails or the response cannot be parsed.
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

	// Parse response header.
	return resp, nil

}
