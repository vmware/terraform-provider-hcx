// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	tflog.Info(ctx, "Creating site pairing")

	client := m.(*Client)

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	tflog.Debug(ctx, "Preparing remote cloud config body", map[string]interface{}{
		"url":      url,
		"username": username,
	})

	body := RemoteCloudConfigBody{
		Remote: RemoteData{
			Username: username,
			Password: password,
			URL:      url,
		},
	}

	res, err := InsertSitePairing(client, body)

	if err != nil {
		tflog.Error(ctx, "Failed to create site pairing", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	secondTry := false
	if res.Errors != nil {
		if res.Errors[0].Error == "Login failure" {
			tflog.Error(ctx, "Login failure when creating site pairing", map[string]interface{}{
				"error_text": res.Errors[0].Text,
			})
			return diag.Errorf("%s", res.Errors[0].Text)
		}

		if len(res.Errors[0].Data) > 0 {
			// Try to get certificate
			tflog.Info(ctx, "Certificate needs to be added for site pairing")
			certificateRaw := res.Errors[0].Data[0]
			certificate, ok := certificateRaw["certificate"].(string)

			if ok {
				// Add certificate
				tflog.Debug(ctx, "Adding certificate for site pairing")
				body := InsertCertificateBody{
					Certificate: certificate,
				}
				_, err := InsertCertificate(client, body)
				if err != nil {
					tflog.Error(ctx, "Failed to add certificate", map[string]interface{}{
						"error": err.Error(),
					})
					return diag.FromErr(err)
				}
				tflog.Info(ctx, "Certificate added successfully")
			}
		} else {
			tflog.Error(ctx, "Unknown error during site pairing", map[string]interface{}{
				"errors": fmt.Sprintf("%+v", res.Errors),
			})
			return diag.Errorf("Unknown error(s): %+v", res.Errors)
		}

		secondTry = true
	}

	if secondTry {
		tflog.Info(ctx, "Retrying site pairing after adding certificate")
		res, err = InsertSitePairing(client, body)
		if err != nil {
			tflog.Error(ctx, "Failed to create site pairing on retry", map[string]interface{}{
				"error": err.Error(),
			})
			return diag.FromErr(err)
		}
	}

	// Wait for job completion
	tflog.Info(ctx, "Waiting for site pairing job to complete", map[string]interface{}{
		"job_id": res.Data.JobID,
	})
	count := 0
	for {
		jr, err := GetJobResult(client, res.Data.JobID)
		if err != nil {
			tflog.Error(ctx, "Failed to get job result", map[string]interface{}{
				"error":  err.Error(),
				"job_id": res.Data.JobID,
			})
			return diag.FromErr(err)
		}

		if jr.IsDone {
			tflog.Info(ctx, "Site pairing job completed successfully")
			break
		}

		if jr.DidFail {
			tflog.Error(ctx, "Site pairing job failed", map[string]interface{}{
				"job_id": res.Data.JobID,
			})
			return diag.Errorf("site pairing job failed")
		}

		count++
		if count%6 == 0 { // Log every minute (6 * 10 seconds)
			tflog.Debug(ctx, "Still waiting for site pairing job to complete", map[string]interface{}{
				"job_id":    res.Data.JobID,
				"wait_time": fmt.Sprintf("%d seconds", count*10),
			})
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

	tflog.Info(ctx, "Reading site pairing", map[string]interface{}{
		"id": d.Id(),
	})

	client := m.(*Client)

	url := d.Get("url").(string)

	tflog.Debug(ctx, "Getting site pairings to find matching URL", map[string]interface{}{
		"url": url,
	})
	res, err := GetSitePairings(client)
	if err != nil {
		tflog.Error(ctx, "Failed to get site pairings", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	for _, item := range res.Data.Items {
		if item.URL == url {
			tflog.Info(ctx, "Found matching site pairing", map[string]interface{}{
				"endpoint_id": item.EndpointID,
				"url":         url,
			})
			d.SetId(item.EndpointID)

			tflog.Debug(ctx, "Getting local container info")
			lc, err := GetLocalContainer(client)
			if err != nil {
				tflog.Error(ctx, "Failed to get local container info", map[string]interface{}{
					"error": err.Error(),
				})
				return diag.FromErr(errors.New("cannot get local container info"))
			}

			if err := d.Set("local_vc", lc.VcUUID); err != nil {
				tflog.Error(ctx, "Failed to set local_vc", map[string]interface{}{
					"error": err.Error(),
				})
				return diag.FromErr(err)
			}

			tflog.Debug(ctx, "Getting remote container info")
			rc, err := GetRemoteContainer(client)
			if err != nil {
				tflog.Error(ctx, "Failed to get remote container info", map[string]interface{}{
					"error": err.Error(),
				})
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
	return diags
}

// resourceSitePairingUpdate updates the site pairing configuration.
func resourceSitePairingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceSitePairingRead(ctx, d, m)
}

// resourceSitePairingDelete removes the site pairing configuration and clears the state of the resource in the schema.
func resourceSitePairingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Deleting site pairing", map[string]interface{}{
		"id": d.Id(),
	})

	client := m.(*Client)
	url := d.Get("url").(string)

	tflog.Debug(ctx, "Calling DeleteSitePairings", map[string]interface{}{
		"endpoint_id": d.Id(),
		"url":         url,
	})
	_, err := DeleteSitePairings(client, d.Id())

	if err != nil {
		tflog.Error(ctx, "Failed to delete site pairing", map[string]interface{}{
			"error": err.Error(),
			"id":    d.Id(),
		})
		return diag.FromErr(err)
	}

	// Wait for site pairing deletion
	tflog.Info(ctx, "Waiting for site pairing deletion to complete")
	for {
		res, err := GetSitePairings(client)
		if err != nil {
			tflog.Error(ctx, "Failed to get site pairings while waiting for deletion", map[string]interface{}{
				"error": err.Error(),
			})
			return diag.FromErr(err)
		}

		found := false
		for _, item := range res.Data.Items {
			if item.URL == url {
				found = true
				tflog.Debug(ctx, "Site pairing still exists, continuing to wait")
			}
		}

		if !found {
			tflog.Info(ctx, "Site pairing deletion completed successfully")
			break
		}

		time.Sleep(5 * time.Second)
	}

	d.SetId("")

	return diags
}
