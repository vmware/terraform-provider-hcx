// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"time"

	b64 "encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-hcx/hcx"
)

func resourcevCenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcevCenterCreate,
		ReadContext:   resourcevCenterRead,
		UpdateContext: resourcevCenterUpdate,
		DeleteContext: resourcevCenterDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The URL of the vCenter instance.",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username to authenticate with the vCenter instance.",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password to authenticate with the vCenter instance.",
				Required:    true,
			},
		},
	}
}

func resourcevCenterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	body := hcx.InsertvCenterBody{
		Data: hcx.InsertvCenterData{
			Items: []hcx.InsertvCenterDataItem{
				{
					Config: hcx.InsertvCenterDataItemConfig{
						Username: username,
						Password: b64.StdEncoding.EncodeToString([]byte(password)),
						URL:      url,
					},
				},
			},
		},
	}

	res, err := hcx.InsertvCenter(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.InsertvCenterData.Items[0].Config.UUID)

	// Restart App Daemon
	hcx.AppEngineStop(client)

	// Wait for App Daemon to be stopped
	for {
		jr, err := hcx.GetAppEngineStatus(client)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Result == "STOPPED" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	hcx.AppEngineStart(client)

	// Wait for App Daemon to be started
	for {
		jr, err := hcx.GetAppEngineStatus(client)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Result == "RUNNING" {
			break
		}
		time.Sleep(5 * time.Second)
	}
	// Seems that we need to wait a bit
	time.Sleep(60 * time.Second)

	return resourcevCenterRead(ctx, d, m)
}

func resourcevCenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourcevCenterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourcevCenterRead(ctx, d, m)
}

func resourcevCenterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	hcx.DeletevCenter(client, d.Id())

	return diags
}
