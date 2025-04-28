// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceLocation defines the resource schema for managing the location for an HCX site.
func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,

		Schema: map[string]*schema.Schema{
			"city": {
				Type:        schema.TypeString,
				Description: "The city where the HCX site is located.",
				Optional:    true,
				Default:     "",
			},
			"country": {
				Type:        schema.TypeString,
				Description: "The country where the HCX site is located.",
				Optional:    true,
				Default:     "",
			},
			"cityascii": {
				Type:        schema.TypeString,
				Description: "The city where the HCX site is located.",
				Computed:    true,
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Description: "The latitude coordinate of the HCX site.",
				Optional:    true,
				Default:     0,
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Description: "The longitude coordinate of the HCX site.",
				Optional:    true,
				Default:     0,
			},
			"province": {
				Type:        schema.TypeString,
				Description: "The province where the HCX site is located.",
				Optional:    true,
				Default:     "",
			},
		},
	}
}

// resourceLocationCreate creates the location configuration for an HCX site.
func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceLocationUpdate(ctx, d, m)
}

// resourceLocationRead retrieves the location configuration for an HCX site.
func resourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	resp, err := GetLocation(client)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.City)
	if err := d.Set("cityascii", resp.City); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("country", resp.Country); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("province", resp.Province); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("latitude", resp.Latitude); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("longitude", resp.Longitude); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceLocationUpdate updates the location configuration for an HCX site.
func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	city := d.Get("city").(string)
	country := d.Get("country").(string)
	CityASCII := city
	latitude := d.Get("latitude").(float64)
	longitude := d.Get("longitude").(float64)
	province := d.Get("province").(string)

	body := SetLocationBody{
		City:      city,
		Country:   country,
		CityASCII: CityASCII,
		Latitude:  latitude,
		Longitude: longitude,
		Province:  province,
	}

	err := SetLocation(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(city)

	return resourceLocationRead(ctx, d, m)
}

// resourceLocationDelete removes the location configuration and clears the state of the resource in the schema.
func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	body := SetLocationBody{
		City:      "",
		Country:   "",
		CityASCII: "",
		Latitude:  constants.DefaultLatitude,
		Longitude: constants.DefaultLongitude,
		Province:  "",
	}

	err := SetLocation(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
