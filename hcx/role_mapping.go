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

// RoleMapping represents the association between a role and a list of user groups.
type RoleMapping struct {
	Role       string   `json:"role"`
	UserGroups []string `json:"userGroups"`
}

// RoleMappingResult represents the result of a role mapping operation.
type RoleMappingResult struct {
	IsSuccess      bool   `json:"isSuccess"`
	Message        string `json:"message"`
	HTTPStatusCode int    `json:"httpStatusCode"`
}

// PutRoleMapping sends a PUT request to update role mappings using the provided body and returns the resulting
// RoleMappingResult object. Returns an error if the request fails or the response cannot be parsed.
func PutRoleMapping(c *Client, body []RoleMapping) (RoleMappingResult, error) {

	resp := RoleMappingResult{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s:9443/api/admin/global/config/roleMappings", c.HostURL), &buf)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	// Send the request.
	_, r, err := c.doAdminRequest(req)
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
