// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hcx "github.com/vmware/terraform-provider-hcx/hcx"
)

func resourceSitePairing() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSitePairingCreate,
		ReadContext:   resourceSitePairingRead,
		UpdateContext: resourceSitePairingUpdate,
		DeleteContext: resourceSitePairingDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"local_vc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_endpoint_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSitePairingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	body := hcx.RemoteCloudConfigBody{
		Remote: hcx.Remote_data{
			Username: username,
			Password: password,
			URL:      url,
		},
	}

	res, err := hcx.InsertSitePairing(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	second_try := false
	if res.Errors != nil {
		if res.Errors[0].Error == "Login failure" {
			return diag.Errorf("%s", res.Errors[0].Text)
		}

		if len(res.Errors[0].Data) > 0 {
			// Try to get certificate
			certificate_raw := res.Errors[0].Data[0]
			certificate, ok := certificate_raw["certificate"].(string)

			if ok {
				// Add certificate
				body := hcx.InsertCertificateBody{
					Certificate: certificate,
				}
				_, err := hcx.InsertCertificate(client, body)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		} else {
			return diag.Errorf("Unknown error(s): %+v", res.Errors)
		}

		second_try = true
	}

	if second_try {
		res, err = hcx.InsertSitePairing(client, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Wait for job completion
	count := 0
	for {
		jr, err := hcx.GetJobResult(client, res.Data.JobID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.IsDone {
			break
		}

		if jr.DidFail {
			return diag.Errorf("Site pairing Job did failed")
		}
		time.Sleep(10 * time.Second)
		count = count + 1
		if count > 5 {
			break
		}
	}

	if count > 5 {
		res, err = hcx.InsertSitePairing(client, body)
		if err != nil {
			return diag.FromErr(err)
		}

		// Wait for job completion
		count = 0
		for {
			jr, err := hcx.GetJobResult(client, res.Data.JobID)
			if err != nil {
				return diag.FromErr(err)
			}

			if jr.IsDone {
				break
			}

			if jr.DidFail {
				return diag.Errorf("Site pairing Job did failed")
			}
			time.Sleep(10 * time.Second)
			count = count + 1
			if count > 5 {
				break
			}
		}
	}

	d.SetId(res.Data.JobID)

	return resourceSitePairingRead(ctx, d, m)
}

func resourceSitePairingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	url := d.Get("url").(string)

	res, err := hcx.GetSitePairings(client)

	for _, item := range res.Data.Items {
		if item.URL == url {
			d.SetId(item.EndpointId)

			lc, err := hcx.GetLocalContainer(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get local container info"))
			}

			d.Set("local_vc", lc.Vcuuid)

			rc, err := hcx.GetRemoteContainer(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get remote container info"))
			}
			d.Set("remote_resource_id", rc.ResourceId)
			d.Set("remote_resource_type", rc.ResourceType)
			d.Set("remote_resource_name", rc.ResourceName)

			// Update Remote Cloud Info
			res2, err := hcx.GetRemoteCloudList(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get remote cloud info"))
			}
			for _, j := range res2.Data.Items {
				if j.URL == url {
					d.Set("remote_name", j.Name)
					d.Set("remote_endpoint_type", res2.Data.Items[0].EndpointType)
				}
			}

			// Update Local Cloud Info
			res3, err := hcx.GetLocalCloudList(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get remote cloud info"))
			}
			d.Set("local_endpoint_id", res3.Data.Items[0].EndpointId)
			d.Set("local_name", res3.Data.Items[0].Name)

			return diags
		}
	}
	if err != nil {
		return diag.FromErr(errors.New("cannot find site pairing info"))
	}

	return diags
}

func resourceSitePairingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceSitePairingRead(ctx, d, m)
}

func resourceSitePairingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)
	url := d.Get("url").(string)

	_, err := hcx.DeleteSitePairings(client, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for site pairing deletion
	for {
		res, err := hcx.GetSitePairings(client)
		if err != nil {
			return diag.FromErr(err)
		}

		found := false
		for _, item := range res.Data.Items {
			if item.URL == url {
				found = true
			}
		}

		if !found {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return diags
}
