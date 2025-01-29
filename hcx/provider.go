// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the schema.Provider object configured with resources, data sources, and schema for the HCX provider.
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
			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_ALLOW_UNVERIFIED_SSL", false),
				Description: "Allow SSL connections with unverifiable certificates.",
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
			"hcx_activation":      resourceActivation(),
			"hcx_compute_profile": resourceComputeProfile(),
			"hcx_l2_extension":    resourceL2Extension(),
			"hcx_location":        resourceLocation(),
			"hcx_network_profile": resourceNetworkProfile(),
			"hcx_rolemapping":     resourceRoleMapping(),
			"hcx_service_mesh":    resourceServiceMesh(),
			"hcx_site_pairing":    resourceSitePairing(),
			"hcx_sso":             resourceSSO(),
			"hcx_vcenter":         resourcevCenter(),
			"hcx_vmc":             resourceVmc(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hcx_compute_profile": dataSourceComputeProfile(),
			"hcx_network_backing": dataSourceNetworkBacking(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	hcxURL := d.Get("hcx").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	adminUsername := d.Get("admin_username").(string)
	adminPassword := d.Get("admin_password").(string)
	allowUnverifiedSSL := d.Get("allow_unverified_ssl").(bool)
	vmcToken := d.Get("vmc_token").(string)

	c, err := NewClient(&hcxURL, &username, &password, &adminUsername, &adminPassword, &allowUnverifiedSSL, &vmcToken)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	if hcxURL == "" {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "No HCX URL provided.",
			Detail:        "Only 'hcx_vmc' resource will be manageable.",
			AttributePath: cty.Path{cty.GetAttrStep{Name: "hcx"}},
		})
	}

	c.Token = vmcToken

	return c, diags
}
