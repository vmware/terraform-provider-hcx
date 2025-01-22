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

// Client -
type Client struct {
	HostURL         string
	HTTPClient      *http.Client
	Token           string
	HcxToken        string
	AdminUsername   string
	AdminPassword   string
	Username        string
	Password        string
	IsAuthenticated bool
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Content struct {
	Strings []string `xml:"string"`
}

type Entries struct {
	Entry []Content `xml:"entry"`
}

// HCX Authentication
func (c *Client) HcxConnectorAuthenticate() error {

	rb, err := json.Marshal(AuthStruct{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		return err
	}

	// authenticate
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/sessions", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var resp *http.Response
	for {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			// Hum...let's wait a bit and try again
			time.Sleep(180 * time.Second)
			resp, err = c.HTTPClient.Do(req)

			if err != nil {
				return fmt.Errorf("Unable to authenticate. Check vCenter User / SSO configuration. Error: %s", err.Error())
			}

		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusAccepted {
			break
		}

		// Check if SSO is ready
		var xmlmessage Entries
		err = xml.Unmarshal(body, &xmlmessage)
		if err != nil {
			return err
		}

		certificate_pb := false
		for _, j := range xmlmessage.Entry {
			if j.Strings[0] == "message" {
				if j.Strings[1] == "'Trusted root certificates' value should not be empty" {
					certificate_pb = true
					log.Println("Certificate error")
				}
			}
		}

		if !certificate_pb {
			return fmt.Errorf("body: %s", body)
		}

		time.Sleep(10 * time.Second)

	}

	// parse response header
	c.Token = resp.Header.Get("x-hm-authorization")

	return nil

}

// NewClient -
func NewClient(hcx, username *string, password *string, adminusername *string, adminpassword *string, vmc_token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		HostURL:         *hcx,
		Username:        *username,
		Password:        *password,
		AdminUsername:   *adminusername,
		AdminPassword:   *adminpassword,
		IsAuthenticated: false,
		Token:           *vmc_token,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, []byte, error) {

	if !c.IsAuthenticated {
		err := c.HcxConnectorAuthenticate()

		if err != nil {
			return nil, nil, err
		}
		c.IsAuthenticated = true
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hm-authorization", fmt.Sprintf("%s", c.Token))

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusAccepted {
			return nil, nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
		}
	}

	return res, body, err
}

func (c *Client) doAdminRequest(req *http.Request) (*http.Response, []byte, error) {

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	c.HTTPClient.Timeout = 300 * time.Second

	if (c.AdminUsername == "") || (c.AdminPassword == "") {
		return nil, nil, fmt.Errorf("Admin Username or Password Empty")
	}

	req.SetBasicAuth(c.AdminUsername, c.AdminPassword)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusNoContent {
			if res.StatusCode != http.StatusAccepted {
				return nil, nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
			}
		}
	}

	return res, body, err
}

func (c *Client) doVmcRequest(req *http.Request) (*http.Response, []byte, error) {

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.HcxToken != "" {
		req.Header.Set("x-hm-authorization", fmt.Sprintf("%s", c.HcxToken))
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusAccepted {
			return nil, nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
		}
	}

	return res, body, err
}
