// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hcx "github.com/vmware/terraform-provider-hcx/hcx"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hcx": {
				Type:        schema.TypeString,
				Description: "The URL of the HCX connector",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username to authenticate for HCX consumption.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password to authenticate for HCX consumption.",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_PASSWORD", nil),
			},
			"admin_username": {
				Type:        schema.TypeString,
				Description: "The username to authenticate with the HCX appliance.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_ADMIN_USER", nil),
			},
			"admin_password": {
				Type:        schema.TypeString,
				Description: "The password to authenticate with the HCX connector",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_ADMIN_PASSWORD", nil),
			},
			"vmc_token": {
				Type:        schema.TypeString,
				Description: "The token to authenticate with the VMware Cloud Services API.",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("VMC_API_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hcx_site_pairing":    resourceSitePairing(),
			"hcx_network_profile": resourceNetworkProfile(),
			"hcx_compute_profile": resourceComputeProfile(),
			"hcx_service_mesh":    resourceServiceMesh(),
			"hcx_l2_extension":    resourceL2Extension(),
			"hcx_vcenter":         resourcevCenter(),
			"hcx_sso":             resourceSSO(),
			"hcx_activation":      resourceActivation(),
			"hcx_rolemapping":     resourceRoleMapping(),
			"hcx_location":        resourceLocation(),
			"hcx_vmc":             resourceVmc(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hcx_network_backing": dataSourceNetworkBacking(),
			"hcx_compute_profile": dataSourceComputeProfile(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	hcxurl := d.Get("hcx").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	adminusername := d.Get("admin_username").(string)
	adminpassword := d.Get("admin_password").(string)
	vmc_token := d.Get("vmc_token").(string)

	c, err := hcx.NewClient(&hcxurl, &username, &password, &adminusername, &adminpassword, &vmc_token)
	//c := &http.Client{Timeout: 10 * time.Second}

	if err != nil {
		return nil, diag.FromErr(err)
	}

	if hcxurl == "" {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "No HCX Url provided",
			Detail:        "Only hcx_vmc resource will be manageable",
			AttributePath: cty.Path{cty.GetAttrStep{Name: "hcx"}},
		})
	}

	c.Token = vmc_token

	return c, diags
}
