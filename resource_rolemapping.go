// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-hcx/hcx"
)

// resourceRoleMapping defines the resource schema for managing role mapping configuration.
func resourceRoleMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleMappingCreate,
		ReadContext:   resourceRoleMappingRead,
		UpdateContext: resourceRoleMappingUpdate,
		DeleteContext: resourceRoleMappingDelete,

		Schema: map[string]*schema.Schema{
			"admin": {
				Type:        schema.TypeList,
				Description: "The group for 'admin' users.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_group": {
							Type:        schema.TypeString,
							Description: "The group name.",
							Optional:    true,
							Default:     "vsphere.local\\Administrators",
						},
					},
				},
			},
			"enterprise": {
				Type:        schema.TypeList,
				Description: "The group for 'enterprise' users.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_group": {
							Type:        schema.TypeString,
							Description: "The group name.",
							Optional:    true,
							Default:     "",
						},
					},
				},
			},
			"sso": {
				Type:        schema.TypeString,
				Description: "The ID of the SSO Lookup Service.",
				Required:    true,
			},
		},
	}
}

// resourceRoleMappingCreate creates the role mapping configuration.
func resourceRoleMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceRoleMappingUpdate(ctx, d, m)
}

// resourceRoleMappingRead retrieves the role mapping configuration.
func resourceRoleMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

// resourceRoleMappingUpdate updates the role mapping configuration.
func resourceRoleMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	admin := d.Get("admin").([]interface{})
	enterprise := d.Get("enterprise").([]interface{})

	adminGroups := []string{}
	for _, j := range admin {
		tmp := j.(map[string]interface{})
		adminGroups = append(adminGroups, tmp["user_group"].(string))
	}

	enterpriseGroups := []string{}
	for _, j := range enterprise {
		tmp := j.(map[string]interface{})
		enterpriseGroups = append(enterpriseGroups, tmp["user_group"].(string))
	}

	body := []hcx.RoleMapping{
		{
			Role:       "System Administrator",
			UserGroups: adminGroups,
		},
		{
			Role:       "Enterprise Administrator",
			UserGroups: enterpriseGroups,
		},
	}
	/*
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		return diag.Errorf("%s", buf)
	*/
	_, err := hcx.PutRoleMapping(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("role_mapping")

	return resourceRoleMappingRead(ctx, d, m)
}

// resourceRoleMappingDelete removes the role mapping configuration amd clears the state of the resource in the schema.
func resourceRoleMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)
	body := []hcx.RoleMapping{
		{
			Role:       "System Administrator",
			UserGroups: []string{},
		},
		{
			Role:       "Enterprise Administrator",
			UserGroups: []string{},
		},
	}

	_, err := hcx.PutRoleMapping(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
