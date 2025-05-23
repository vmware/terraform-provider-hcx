// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSitePairing defines the resource schema for managing site pairing configuration.
func resourceSitePairing() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSitePairingCreate,
		ReadContext:   resourceSitePairingRead,
		UpdateContext: resourceSitePairingUpdate,
		DeleteContext: resourceSitePairingDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The URL of the remote cloud.",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username used for remote cloud authentication.",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password used for remote cloud authentication.",
				Required:    true,
				Sensitive:   true,
			},
			"local_vc": {
				Type:        schema.TypeString,
				Description: "The ID of the local vCenter instance.",
				Computed:    true,
			},
			"local_endpoint_id": {
				Type:        schema.TypeString,
				Description: "The endpoint ID of the local HCX site.",
				Computed:    true,
			},
			"local_name": {
				Type:        schema.TypeString,
				Description: "The endpoint name of the local HCX site.",
				Computed:    true,
			},
			"remote_name": {
				Type:        schema.TypeString,
				Description: "The endpoint name of the remote HCX site.",
				Computed:    true,
			},
			"remote_endpoint_type": {
				Type:        schema.TypeString,
				Description: "The endpoint type of the remote HCX site.",
				Computed:    true,
			},
			"remote_resource_id": {
				Type:        schema.TypeString,
				Description: "The resource ID of the remote cloud.",
				Computed:    true,
			},
			"remote_resource_name": {
				Type:        schema.TypeString,
				Description: "The resource name of the remote HCX site.",
				Computed:    true,
			},
			"remote_resource_type": {
				Type:        schema.TypeString,
				Description: "The resource type of the remote HCX site.",
				Computed:    true,
			},
		},
	}
}

// resourceSitePairingCreate creates the site paring configuration.
func resourceSitePairingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	body := RemoteCloudConfigBody{
		Remote: RemoteData{
			Username: username,
			Password: password,
			URL:      url,
		},
	}

	res, err := InsertSitePairing(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	secondTry := false
	if res.Errors != nil {
		if res.Errors[0].Error == "Login failure" {
			return diag.Errorf("%s", res.Errors[0].Text)
		}

		if len(res.Errors[0].Data) > 0 {
			// Try to get certificate
			certificateRaw := res.Errors[0].Data[0]
			certificate, ok := certificateRaw["certificate"].(string)

			if ok {
				// Add certificate
				body := InsertCertificateBody{
					Certificate: certificate,
				}
				_, err := InsertCertificate(client, body)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		} else {
			return diag.Errorf("Unknown error(s): %+v", res.Errors)
		}

		secondTry = true
	}

	if secondTry {
		res, err = InsertSitePairing(client, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Wait for job completion
	count := 0
	for {
		jr, err := GetJobResult(client, res.Data.JobID)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.IsDone {
			break
		}

		if jr.DidFail {
			return diag.Errorf("site pairing job failed")
		}
		time.Sleep(10 * time.Second)
		count = count + 1
		if count > 5 {
			break
		}
	}

	if count > 5 {
		res, err = InsertSitePairing(client, body)
		if err != nil {
			return diag.FromErr(err)
		}

		// Wait for job completion
		count = 0
		for {
			jr, err := GetJobResult(client, res.Data.JobID)
			if err != nil {
				return diag.FromErr(err)
			}

			if jr.IsDone {
				break
			}

			if jr.DidFail {
				return diag.Errorf("site pairing job failed")
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

// resourceSitePairingRead retrieves a site paring configuration.
func resourceSitePairingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	url := d.Get("url").(string)

	res, err := GetSitePairings(client)

	for _, item := range res.Data.Items {
		if item.URL == url {
			d.SetId(item.EndpointID)

			lc, err := GetLocalContainer(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get local container info"))
			}

			if err := d.Set("local_vc", lc.VcUUID); err != nil {
				return diag.FromErr(err)
			}

			rc, err := GetRemoteContainer(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get remote container info"))
			}
			if err := d.Set("remote_resource_id", rc.ResourceID); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("remote_resource_type", rc.ResourceType); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("remote_resource_name", rc.ResourceName); err != nil {
				return diag.FromErr(err)
			}

			// Update Remote Cloud Info
			res2, err := GetRemoteCloudList(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get remote cloud info"))
			}
			for _, j := range res2.Data.Items {
				if j.URL == url {
					if err := d.Set("remote_name", j.Name); err != nil {
						return diag.FromErr(err)
					}
					if err := d.Set("remote_endpoint_type", res2.Data.Items[0].EndpointType); err != nil {
						return diag.FromErr(err)
					}
				}
			}

			// Update Local Cloud Info
			res3, err := GetLocalCloudList(client)
			if err != nil {
				return diag.FromErr(errors.New("cannot get remote cloud info"))
			}
			if err := d.Set("local_endpoint_id", res3.Data.Items[0].EndpointID); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("local_name", res3.Data.Items[0].Name); err != nil {
				return diag.FromErr(err)
			}

			return diags
		}
	}
	if err != nil {
		return diag.FromErr(errors.New("cannot find site pairing info"))
	}

	return diags
}

// resourceSitePairingUpdate updates the site pairing configuration.
func resourceSitePairingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceSitePairingRead(ctx, d, m)
}

// resourceSitePairingDelete removes the site pairing configuration and clears the state of the resource in the schema.
func resourceSitePairingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)
	url := d.Get("url").(string)

	_, err := DeleteSitePairings(client, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for site pairing deletion
	for {
		res, err := GetSitePairings(client)
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
