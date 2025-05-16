// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"time"

	b64 "encoding/base64"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	tflog.Info(ctx, "Creating vCenter instance configuration", map[string]interface{}{
		"url":      url,
		"username": username,
	})

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
		tflog.Error(ctx, "Failed to create vCenter instance", map[string]interface{}{
			"error":    err.Error(),
			"url":      url,
			"username": username,
		})
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "vCenter instance created successfully", map[string]interface{}{
		"uuid": res.InsertvCenterData.Items[0].Config.UUID,
	})
	d.SetId(res.InsertvCenterData.Items[0].Config.UUID)

	// Restart App Daemon
	tflog.Info(ctx, "Stopping app engine to apply vCenter changes")
	_, err = AppEngineStop(client)
	if err != nil {
		tflog.Error(ctx, "Failed to stop app engine", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	// Wait for App Daemon to be stopped
	tflog.Debug(ctx, "Waiting for app engine to stop")
	for {
		jr, err := GetAppEngineStatus(client)
		if err != nil {
			tflog.Error(ctx, "Failed to get app engine status", map[string]interface{}{
				"error": err.Error(),
			})
			return diag.FromErr(err)
		}

		if jr.Result == constants.StoppedStatus {
			tflog.Debug(ctx, "App engine stopped successfully")
			break
		}
		tflog.Debug(ctx, "App engine stopping, waiting...", map[string]interface{}{
			"status": jr.Result,
		})
		time.Sleep(5 * time.Second)
	}

	tflog.Info(ctx, "Starting app engine after vCenter configuration")
	_, err = AppEngineStart(client)
	if err != nil {
		tflog.Error(ctx, "Failed to start app engine", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	// Wait for App Daemon to be started
	tflog.Debug(ctx, "Waiting for app engine to start")
	for {
		jr, err := GetAppEngineStatus(client)
		if err != nil {
			tflog.Error(ctx, "Failed to get app engine status", map[string]interface{}{
				"error": err.Error(),
			})
			return diag.FromErr(err)
		}

		if jr.Result == constants.RunningStatus {
			tflog.Debug(ctx, "App engine started successfully")
			break
		}
		tflog.Debug(ctx, "App engine starting, waiting...", map[string]interface{}{
			"status": jr.Result,
		})
		time.Sleep(5 * time.Second)
	}
	// Seems that we need to wait a bit
	tflog.Info(ctx, "Waiting for services to initialize completely")
	time.Sleep(60 * time.Second)

	return resourcevCenterRead(ctx, d, m)
}

// resourcevCenterRead retrieves the vCenter instance configuration.
func resourcevCenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tflog.Debug(ctx, "Reading vCenter instance configuration", map[string]interface{}{
		"uuid": d.Id(),
	})
	return diags
}

// resourcevCenterUpdate updates the vCenter instance configuration.
func resourcevCenterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Updating vCenter instance configuration", map[string]interface{}{
		"uuid": d.Id(),
	})
	return resourcevCenterRead(ctx, d, m)
}

// resourcevCenterDelete removes the vCenter instance and clears the state of the resource in the schema.
func resourcevCenterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Client)

	tflog.Info(ctx, "Deleting vCenter instance configuration", map[string]interface{}{
		"uuid": d.Id(),
	})

	_, err := DeletevCenter(client, d.Id())
	if err != nil {
		tflog.Error(ctx, "Failed to delete vCenter instance", map[string]interface{}{
			"error": err.Error(),
			"uuid":  d.Id(),
		})
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "vCenter instance deleted successfully")
	return diags
}
