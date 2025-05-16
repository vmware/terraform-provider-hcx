// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
func (c *Client) HcxConnectorAuthenticate(ctx context.Context) error {
	tflog.Info(ctx, "Starting HCX connector authentication", map[string]interface{}{
		"host_url": c.HostURL,
	})

	rb, err := json.Marshal(AuthStruct{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		tflog.Error(ctx, "Failed to marshal authentication request body", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to marshal authentication request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/sessions", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		tflog.Error(ctx, "Failed to create authentication request", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to create authentication request: %w", err)
	}

	req = req.WithContext(ctx)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	var resp *http.Response
	for {
		tflog.Debug(ctx, "Attempting authentication request")
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			tflog.Warn(ctx, "Authentication attempt failed, will retry after delay", map[string]interface{}{
				"error": err.Error(),
			})
			time.Sleep(180 * time.Second)
			resp, err = c.HTTPClient.Do(req)

			if err != nil {
				tflog.Error(ctx, "Authentication failed after retry", map[string]interface{}{
					"error": err.Error(),
				})
				return fmt.Errorf("authentication failed after retry; check credentials: %w", err)
			}
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			tflog.Error(ctx, "Failed to read authentication response body", map[string]interface{}{
				"error": err.Error(),
			})
			return fmt.Errorf("failed to read authentication response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			tflog.Info(ctx, "Authentication successful", map[string]interface{}{
				"status_code": resp.StatusCode,
			})
			break
		}

		tflog.Debug(ctx, "Checking if SSO is ready", map[string]interface{}{
			"status_code": resp.StatusCode,
			"body_length": len(body),
		})

		// Check if SSO is ready.
		var xmlmessage Entries
		err = xml.Unmarshal(body, &xmlmessage)
		if err != nil {
			tflog.Error(ctx, "Failed to unmarshal XML response", map[string]interface{}{
				"error": err.Error(),
				"body":  string(body),
			})
			return fmt.Errorf("failed to unmarshal XML response: %w", err)
		}

		certificatePb := false
		for _, j := range xmlmessage.Entry {
			if j.Strings[0] == "message" {
				if j.Strings[1] == "'Trusted root certificates' value should not be empty" {
					certificatePb = true
					tflog.Info(ctx, "Certificate error detected, will retry")
				}
			}
		}

		if !certificatePb {
			tflog.Error(ctx, "Unexpected authentication response", map[string]interface{}{
				"body": string(body),
			})
			return fmt.Errorf("unexpected authentication response body: %s", body)
		}

		tflog.Info(ctx, "Waiting to retry authentication due to certificate issue")
		time.Sleep(10 * time.Second)
	}

	// Parse response header.
	c.Token = resp.Header.Get("x-hm-authorization")
	tflog.Debug(ctx, "Authentication token retrieved successfully")

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
	ctx := req.Context()

	if !c.IsAuthenticated {
		tflog.Debug(ctx, "Client not authenticated, performing authentication")
		err := c.HcxConnectorAuthenticate(ctx)
		if err != nil {
			tflog.Error(ctx, "Authentication failed during request", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, nil, fmt.Errorf("authentication failed during request: %w", err)
		}
		c.IsAuthenticated = true
		tflog.Info(ctx, "Authentication successful")
	}

	tflog.Debug(ctx, "Sending request", map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	})

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hm-authorization", c.Token)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Error(ctx, "Failed to read HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		tflog.Warn(ctx, "Unexpected response status", map[string]interface{}{
			"status_code": res.StatusCode,
			"body":        string(body),
		})
		return nil, nil, fmt.Errorf("unexpected response status: %d, body: %s", res.StatusCode, body)
	}

	tflog.Debug(ctx, "Request successful", map[string]interface{}{
		"status_code": res.StatusCode,
	})

	return res, body, nil
}

// doAdminRequest executes an HTTP request using the admin credentials for Basic Authentication. It supports requests
// that require elevated permissions and optionally skips SSL verification. Returns the response, response body, and any
// encountered error.
func (c *Client) doAdminRequest(req *http.Request) (*http.Response, []byte, error) {
	ctx := req.Context()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	c.HTTPClient.Timeout = 300 * time.Second

	if (c.AdminUsername == "") || (c.AdminPassword == "") {
		tflog.Error(ctx, "Admin credentials missing")
		return nil, nil, fmt.Errorf("admin_username or admin_password is empty")
	}

	req.SetBasicAuth(c.AdminUsername, c.AdminPassword)

	tflog.Debug(ctx, "Sending admin request", map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	})

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send admin HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Error(ctx, "Failed to read admin HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	tflog.Debug(ctx, "Received admin response", map[string]interface{}{
		"status_code": res.StatusCode,
		"body_length": len(body),
	})

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusAccepted {
		tflog.Error(ctx, "Unexpected admin response status", map[string]interface{}{
			"status_code": res.StatusCode,
			"body":        string(body),
		})
		return nil, nil, fmt.Errorf("unexpected response status from admin request: %d, body: %s", res.StatusCode, body)
	}

	return res, body, nil
}

// doVmcRequest sends an HTTP request to the VMware Cloud (VMC) service. It uses the HCX token for authorization if
// present and optionally skips SSL verification. Returns the response, response body, and any encountered error.
func (c *Client) doVmcRequest(req *http.Request) (*http.Response, []byte, error) {
	ctx := req.Context()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.HcxToken != "" {
		req.Header.Set("x-hm-authorization", c.HcxToken)
	} else {
		tflog.Debug(ctx, "No HCX token available for VMC request")
	}

	tflog.Debug(ctx, "Sending VMC request", map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	})

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.AllowUnverifiedSSL, // #nosec G402
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		tflog.Error(ctx, "Failed to send VMC HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Error(ctx, "Failed to read VMC HTTP response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	tflog.Debug(ctx, "Received VMC response", map[string]interface{}{
		"status_code": res.StatusCode,
		"body_length": len(body),
	})

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		tflog.Error(ctx, "Unexpected VMC response status", map[string]interface{}{
			"status_code": res.StatusCode,
			"body":        string(body),
		})
		return nil, nil, fmt.Errorf("unexpected vmc response status: %d, body: %s", res.StatusCode, body)
	}

	return res, body, nil
}
