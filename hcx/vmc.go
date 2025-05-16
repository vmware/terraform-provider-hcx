// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

// VmcAuthenticate sends a request to authenticate with the VMware Cloud (VMC) API using the provided token.
// Returns an access token as a string or an error if the request fails or the response cannot be parsed.
func VmcAuthenticate(token string) (string, error) {
	ctx := context.Background()
	tflog.Info(ctx, "Authenticating with VMware Cloud API")

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.VmcAuthURL,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/api-tokens/authorize?refresh_token=%s", c.HostURL, token), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create VMC authentication POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send VMC authentication POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to send POST request: %w", err)
	}

	resp := VmcAccessToken{}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse VMC authentication HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Info(ctx, "Successfully authenticated with VMware Cloud API")
	return resp.AccessToken, nil
}

// CloudAuthenticate sends a request to authenticate to the HCX cloud service using the provided token.
// On success, it sets the HcxToken field of the provided Client.
func CloudAuthenticate(client *Client, token string) error {
	ctx := context.Background()
	tflog.Info(ctx, "Authenticating to HCX Cloud service")

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
		tflog.Error(ctx, "Failed to encode HCX Cloud authentication request body", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sessions", c.HostURL), &buf)
	if err != nil {
		tflog.Error(ctx, "Failed to create HCX Cloud authentication POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to create POST request: %w", err)
	}

	resp, _, err := c.doVmcRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send HCX Cloud authentication POST request", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to send POST request: %w", err)
	}

	auth := resp.Header.Get("x-hm-authorization")
	if auth == "" {
		tflog.Error(ctx, "Failed to authorize: x-hm-authorization header not found")
		return errors.New("failed to authorize: x-hm-authorization header not found")
	}

	client.HcxToken = auth
	tflog.Info(ctx, "Successfully authenticated to HCX Cloud service")
	return nil
}

// GetSddcByName sends a request to retrieve an SDDC by name.
// Returns the matching SDDC object or an error.
func GetSddcByName(client *Client, sddcName string) (SDDC, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Retrieving SDDC by name", map[string]interface{}{
		"sddcName": sddcName,
	})

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/sddcs", c.HostURL), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create SDDC GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return SDDC{}, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send SDDC GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return SDDC{}, fmt.Errorf("failed to send GET request: %w", err)
	}

	resp := GetSddcsResults{}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse SDDC HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return SDDC{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	for _, j := range resp.SDDCs {
		if j.Name == sddcName {
			tflog.Debug(ctx, "Found SDDC by name", map[string]interface{}{
				"sddcName": sddcName,
				"sddcID":   j.ID,
			})
			return j, nil
		}
	}

	tflog.Error(ctx, "Failed to find SDDC by name", map[string]interface{}{
		"sddcName": sddcName,
	})
	return SDDC{}, errors.New("failed to find SDDC by name")
}

// GetSddcByID sends a request to retrieve an SDDC by ID.
func GetSddcByID(client *Client, sddcID string) (SDDC, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Retrieving SDDC by ID", map[string]interface{}{
		"sddcID": sddcID,
	})

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/sddcs", c.HostURL), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create SDDC GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return SDDC{}, fmt.Errorf("failed to create GET request: %w", err)
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send SDDC GET request", map[string]interface{}{
			"error": err.Error(),
		})
		return SDDC{}, fmt.Errorf("failed to send GET request: %w", err)
	}

	resp := GetSddcsResults{}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse SDDC HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return SDDC{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	for _, j := range resp.SDDCs {
		if j.ID == sddcID {
			tflog.Debug(ctx, "Found SDDC by ID", map[string]interface{}{
				"sddcID":   sddcID,
				"sddcName": j.Name,
			})
			return j, nil
		}
	}

	tflog.Error(ctx, "Failed to find SDDC by ID", map[string]interface{}{
		"sddcID": sddcID,
	})
	return SDDC{}, errors.New("failed to find SDDC by ID")
}

// ActivateHcxOnSDDC sends a request to activate HCX on the specified SDDC.
func ActivateHcxOnSDDC(client *Client, sddcID string) (ActivateHcxOnSDDCResults, error) {
	ctx := context.Background()
	tflog.Info(ctx, "Activating HCX on SDDC", map[string]interface{}{
		"sddcID": sddcID,
	})

	resp := ActivateHcxOnSDDCResults{}

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sddcs/%s?action=activate", c.HostURL, sddcID), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create HCX activation POST request", map[string]interface{}{
			"error":  err.Error(),
			"sddcID": sddcID,
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send HCX activation POST request", map[string]interface{}{
			"error":  err.Error(),
			"sddcID": sddcID,
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse HCX activation HTTP response", map[string]interface{}{
			"error":  err.Error(),
			"sddcID": sddcID,
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Info(ctx, "HCX activation initiated on SDDC", map[string]interface{}{
		"sddcID": sddcID,
		"jobID":  resp.JobID,
	})
	return resp, nil
}

// DeactivateHcxOnSDDC sends a request to deactivate HCX on the specified SDDC.
func DeactivateHcxOnSDDC(client *Client, sddcID string) (DeactivateHcxOnSDDCResults, error) {
	ctx := context.Background()
	tflog.Info(ctx, "Deactivating HCX on SDDC", map[string]interface{}{
		"sddcID": sddcID,
	})

	resp := DeactivateHcxOnSDDCResults{}

	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    constants.HcxCloudConsumerURL,
		HcxToken:   client.HcxToken,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sddcs/%s?action=deactivate", c.HostURL, sddcID), nil)
	if err != nil {
		tflog.Error(ctx, "Failed to create HCX deactivation POST request", map[string]interface{}{
			"error":  err.Error(),
			"sddcID": sddcID,
		})
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doVmcRequest(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send HCX deactivation POST request", map[string]interface{}{
			"error":  err.Error(),
			"sddcID": sddcID,
		})
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		tflog.Error(ctx, "Failed to parse HCX deactivation HTTP response", map[string]interface{}{
			"error":  err.Error(),
			"sddcID": sddcID,
		})
		return resp, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	tflog.Info(ctx, "HCX deactivation initiated on SDDC", map[string]interface{}{
		"sddcID": sddcID,
		"jobID":  resp.JobID,
	})
	return resp, nil
}
