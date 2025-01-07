package main

import (
	"context"

	"github.com/adeleporte/terraform-provider-hcx/hcx"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,

		Schema: map[string]*schema.Schema{
			"city": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"cityascii": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
			"longitude": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
			"province": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceLocationUpdate(ctx, d, m)
}

func resourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	resp, err := hcx.GetLocation(client)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.City)
	d.Set("cityascii", resp.City)
	d.Set("country", resp.Country)
	d.Set("province", resp.Province)
	d.Set("latitude", resp.Latitude)
	d.Set("longitude", resp.Longitude)

	return diags
}

func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	city := d.Get("city").(string)
	country := d.Get("country").(string)
	cityAscii := city
	latitude := d.Get("latitude").(float64)
	longitude := d.Get("longitude").(float64)
	province := d.Get("province").(string)

	body := hcx.SetLocationBody{
		City:      city,
		Country:   country,
		CityAscii: cityAscii,
		Latitude:  latitude,
		Longitude: longitude,
		Province:  province,
	}

	err := hcx.SetLocation(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(city)

	return resourceLocationRead(ctx, d, m)
}

func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	body := hcx.SetLocationBody{
		City:      "",
		Country:   "",
		CityAscii: "",
		Latitude:  0,
		Longitude: 0,
		Province:  "",
	}

	err := hcx.SetLocation(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
