// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-hcx/hcx"
)

// resourceVmc defines the resource schema for managing a VMware Cloud on AWS configuration.
func resourceVmc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVmcCreate,
		ReadContext:   resourceVmcRead,
		UpdateContext: resourceVmcUpdate,
		DeleteContext: resourceVmcDelete,

		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:        schema.TypeString,
				Description: "The ID of the SDDC.",
				Optional:    true,
			},
			"sddc_name": {
				Type:        schema.TypeString,
				Description: "The name of the SDDC.",
				Optional:    true,
			},
			"cloud_url": {
				Type:        schema.TypeString,
				Description: "The URL of HCX Cloud, used for the site pairing configuration.",
				Computed:    true,
			},
			"cloud_name": {
				Type:        schema.TypeString,
				Description: "The name of the HCX Cloud.",
				Computed:    true,
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Description: "The type of the HCX Cloud. Use 'nsp' for VMware Cloud on AWS.",
				Computed:    true,
			},
		},
	}
}

// resourceVmcCreate creates the VMware Cloud on AWS configuration.
func resourceVmcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcx.Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	if sddcName == "" && sddcID == "" {
		return diag.Errorf("SDDC name or Id must be specified")
	}

	// Authenticate with VMware Cloud Services
	accessToken, err := hcx.VmcAuthenticate(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = hcx.CloudAuthenticate(client, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	var sddc hcx.SDDC
	if sddcID != "" {
		sddc, err = hcx.GetSddcByID(client, sddcID)
	} else {
		sddc, err = hcx.GetSddcByName(client, sddcName)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	// Check if already activated
	if sddc.DeploymentStatus == "ACTIVE" {
		return diag.Errorf("Already activated")
	}

	// Activate HCX
	_, err = hcx.ActivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task to be completed
	errcount := 0
	for {
		if sddcID != "" {
			sddc, err = hcx.GetSddcByID(client, sddcID)
		} else {
			sddc, err = hcx.GetSddcByName(client, sddcName)
		}
		if err != nil {
			// Attempt to bypass recurring situation where the HCX API
			// returns status 502 with a proxy server error, and an HTML response
			// instead of JSON.
			errcount++
			hclog.Default().Info("[INFO] - resourceVmcCreate() - Error retrieving SDDC status: ", "error", err.Error(), "Errcount:", errcount)
			if errcount > 12 {
				return diag.FromErr(err)
			}
		}

		if sddc.DeploymentStatus == "ACTIVE" {
			break
		}

		if sddc.DeploymentStatus == "ACTIVATION_FAILED" {
			return diag.Errorf("Activation failed")
		}

		time.Sleep(10 * time.Second)
	}

	return resourceVmcRead(ctx, d, m)
}

// resourceVmcRead retrieves the VMware Cloud on AWS configuration.
func resourceVmcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	log.Printf("******************************************************************\n")
	log.Printf("token: %s, sddc_name: %s,   sddc: %s   \n", token, sddcName, sddcID)
	log.Printf("******************************************************************\n")

	if sddcName == "" && sddcID == "" {
		return diag.Errorf("SDDC name or Id must be specified")
	}

	// Authenticate with VMware Cloud Services
	accessToken, err := hcx.VmcAuthenticate(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = hcx.CloudAuthenticate(client, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("****************")
	log.Printf("[Client inside]: %+v ", client)
	log.Printf("****************")

	var sddc hcx.SDDC
	if sddcID != "" {
		sddc, err = hcx.GetSddcByID(client, sddcID)
	} else {
		sddc, err = hcx.GetSddcByName(client, sddcName)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sddc.ID)
	d.Set("cloud_url", sddc.CloudURL)
	d.Set("cloud_name", sddc.CloudName)
	d.Set("cloud_type", sddc.CloudType)

	return diags
}

// resourceVmcUpdate updates the VMware Cloud on AWS resource configuration.
func resourceVmcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceVmcRead(ctx, d, m)
}

// resourceVmcDelete removes the VMware Cloud on AWS configuration and clears the state of the resource in the schema.
func resourceVmcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcx.Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	// Authenticate with VMware Cloud Services
	accessToken, err := hcx.VmcAuthenticate(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = hcx.CloudAuthenticate(client, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	var sddc hcx.SDDC
	if sddcID != "" {
		sddc, err = hcx.GetSddcByID(client, sddcID)
	} else {
		sddc, err = hcx.GetSddcByName(client, sddcName)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// Deactivate HCX
	_, err = hcx.DeactivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task to be completed
	errcount := 0
	for {
		var sddc hcx.SDDC
		if sddcID != "" {
			sddc, err = hcx.GetSddcByID(client, sddcID)
		} else {
			sddc, err = hcx.GetSddcByName(client, sddcName)
		}
		if err != nil {
			// Attempt to bypass recurring situation where the HCX API
			// returns status 502 with a proxy server error, and an HTML response
			// instead of JSON.
			errcount++
			hclog.Default().Info("[INFO] - resourceVmcDelete() - Error retrieving SDDC status: ", "error", err.Error(), "Errcount:", errcount)
			if errcount > 12 {
				return diag.FromErr(err)
			}
		}

		if sddc.DeploymentStatus == "DE-ACTIVATED" {
			break
		}

		if sddc.DeploymentStatus == "DEACTIVATION_FAILED" {
			return diag.Errorf("Deactivation failed")
		}

		time.Sleep(10 * time.Second)
	}

	return diags
}
