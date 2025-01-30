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

// InsertCertificateBody represents the body structure used to insert a certificate.
type InsertCertificateBody struct {
	Certificate string `json:"certificate"`
}

// InsertCertificateResult represents the result of inserting a certificate, including success and completion status.
type InsertCertificateResult struct {
	Success   bool `json:"success"`
	Completed bool `json:"completed"`
}

// InsertCertificate sends a request to create a new certificate using the provided body and returns an
// InsertCertificateResult object. Returns an error if the request fails or the response cannot be parsed.
func InsertCertificate(c *Client, body InsertCertificateBody) (InsertCertificateResult, error) {

	resp := InsertCertificateResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/admin/certificates", c.HostURL), &buf)
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
