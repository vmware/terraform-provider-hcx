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

type InsertCertificateBody struct {
	Certificate string `json:"certificate"`
}

type InsertCertificateResult struct {
	Success   bool `json:"success"`
	Completed bool `json:"completed"`
}

// InsertL2Extention ...
func InsertCertificate(c *Client, body InsertCertificateBody) (InsertCertificateResult, error) {

	resp := InsertCertificateResult{}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hybridity/api/admin/certificates", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request
	_, r, err := c.doRequest(req)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// parse response body
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	return resp, nil
}
