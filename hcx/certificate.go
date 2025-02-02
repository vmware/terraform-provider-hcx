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
		return resp, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/admin/certificates", c.HostURL), &buf)
	if err != nil {
		return resp, fmt.Errorf("failed to create POST request: %w", err)
	}

	_, r, err := c.doRequest(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send POST request: %w", err)
	}

	err = json.Unmarshal(r, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal POST response: %w", err)
	}

	return resp, nil
}
