// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client represents a structure for managing HTTP communication and authentication details.
type Client struct {
	HostURL            string
	HTTPClient         *http.Client
	Token              string
	HcxToken           string
	AdminUsername      string
	AdminPassword      string
	Username           string
	Password           string
	IsAuthenticated    bool
	AllowUnverifiedSSL bool
}

// AuthStruct represents a structure containing username and password for authentication purposes.
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the structure of a response returned after a successful authentication.
type AuthResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// Content represents a reusable structure containing a slice of strings, typically used for XML unmarshaling.
type Content struct {
	Strings []string `xml:"string"`
}

// Entries represents a collection of content entries, typically used for XML unmarshaling.
type Entries struct {
	Entry []Content `xml:"entry"`
}

// HcxConnectorAuthenticate authenticates the client with the HCX service by sending a request with user credentials.
// It retrieves and stores the HCX authorization token required for subsequent requests.
func (c *Client) HcxConnectorAuthenticate() error {

	rb, err := json.Marshal(AuthStruct{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal authentication request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/sessions", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("failed to create authentication request: %w", err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	var resp *http.Response
	for {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			time.Sleep(180 * time.Second)
			resp, err = c.HTTPClient.Do(req)

			if err != nil {
				return fmt.Errorf("authentication failed after retry; check credentials: %w", err)
			}
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read authentication response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			break
		}

		// Check if SSO is ready.
		var xmlmessage Entries
		err = xml.Unmarshal(body, &xmlmessage)
		if err != nil {
			return fmt.Errorf("failed to unmarshal XML response: %w", err)
		}

		certificatePb := false
		for _, j := range xmlmessage.Entry {
			if j.Strings[0] == "message" {
				if j.Strings[1] == "'Trusted root certificates' value should not be empty" {
					certificatePb = true
					log.Println("Certificate error")
				}
			}
		}

		if !certificatePb {
			return fmt.Errorf("unexpected authentication response body: %s", body)
		}

		time.Sleep(10 * time.Second)
	}

	// Parse response header.
	c.Token = resp.Header.Get("x-hm-authorization")

	return nil
}

// NewClient initializes and returns a new Client instance with the provided configuration, including authentication
// details, HCX URL, and SSL settings.
func NewClient(hcx, username *string, password *string, adminUsername *string, adminPassword *string, allowUnverifiedSSL *bool, vmcToken *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		HostURL:            *hcx,
		Username:           *username,
		Password:           *password,
		AdminUsername:      *adminUsername,
		AdminPassword:      *adminPassword,
		IsAuthenticated:    false,
		AllowUnverifiedSSL: *allowUnverifiedSSL,
		Token:              *vmcToken,
	}

	return &c, nil
}

// doRequest performs an authenticated HTTP request. If the client is not yet authenticated, it performs authentication
// first, then executes the request. Returns the HTTP response, response body, and any encountered error.
func (c *Client) doRequest(req *http.Request) (*http.Response, []byte, error) {

	if !c.IsAuthenticated {
		err := c.HcxConnectorAuthenticate()
		if err != nil {
			return nil, nil, fmt.Errorf("authentication failed during request: %w", err)
		}
		c.IsAuthenticated = true
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hm-authorization", c.Token)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return nil, nil, fmt.Errorf("unexpected response status: %d, body: %s", res.StatusCode, body)
	}

	return res, body, nil
}

// doAdminRequest executes an HTTP request using the admin credentials for Basic Authentication. It supports requests
// that require elevated permissions and optionally skips SSL verification. Returns the response, response body, and any
// encountered error.
func (c *Client) doAdminRequest(req *http.Request) (*http.Response, []byte, error) {

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	c.HTTPClient.Timeout = 300 * time.Second

	if (c.AdminUsername == "") || (c.AdminPassword == "") {
		return nil, nil, fmt.Errorf("admin_username or admin_password is empty")
	}

	req.SetBasicAuth(c.AdminUsername, c.AdminPassword)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusAccepted {
		return nil, nil, fmt.Errorf("unexpected response status from admin request: %d, body: %s", res.StatusCode, body)
	}

	return res, body, nil
}

// doVmcRequest sends an HTTP request to the VMware Cloud (VMC) service. It uses the HCX token for authorization if
// present and optionally skips SSL verification. Returns the response, response body, and any encountered error.
func (c *Client) doVmcRequest(req *http.Request) (*http.Response, []byte, error) {

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.HcxToken != "" {
		req.Header.Set("x-hm-authorization", c.HcxToken)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return nil, nil, fmt.Errorf("unexpected vmc response status: %d, body: %s", res.StatusCode, body)
	}

	return res, body, nil
}
