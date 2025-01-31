// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"time"

	b64 "encoding/base64"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcevCenter defines the resource schema for managing vCenter instance configuration.
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

// resourcevCenterCreate creates the vCenter instance configuration.
func resourcevCenterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Client)

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	body := InsertvCenterBody{
		Data: InsertvCenterData{
			Items: []InsertvCenterDataItem{
				{
					Config: InsertvCenterDataItemConfig{
						Username: username,
						Password: b64.StdEncoding.EncodeToString([]byte(password)),
						URL:      url,
					},
				},
			},
		},
	}

	res, err := InsertvCenter(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.InsertvCenterData.Items[0].Config.UUID)

	// Restart App Daemon
	_, err = AppEngineStop(client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for App Daemon to be stopped
	for {
		jr, err := GetAppEngineStatus(client)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Result == constants.StoppedStatus {
			break
		}
		time.Sleep(5 * time.Second)
	}

	_, err = AppEngineStart(client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for App Daemon to be started
	for {
		jr, err := GetAppEngineStatus(client)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Result == constants.RunningStatus {
			break
		}
		time.Sleep(5 * time.Second)
	}
	// Seems that we need to wait a bit
	time.Sleep(60 * time.Second)

	return resourcevCenterRead(ctx, d, m)
}

// resourcevCenterRead retrieves the vCenter instance configuration.
func resourcevCenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

// resourcevCenterUpdate updates the vCenter instance configuration.
func resourcevCenterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourcevCenterRead(ctx, d, m)
}

// resourcevCenterDelete removes the vCenter instance and clears the state of the resource in the schema.
func resourcevCenterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

	_, err := DeletevCenter(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
