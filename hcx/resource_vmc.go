// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package hcx

import (
	"context"
	"log"
	"time"

	"github.com/vmware/terraform-provider-hcx/hcx/constants"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:         schema.TypeString,
				Description:  "The ID of the SDDC.",
				Optional:     true,
				ExactlyOneOf: []string{"sddc_id", "sddc_name"}, // Enforces that at least one of them is provided.
			},
			"sddc_name": {
				Type:         schema.TypeString,
				Description:  "The name of the SDDC.",
				Optional:     true,
				ExactlyOneOf: []string{"sddc_id", "sddc_name"}, // Enforces that at least one of them is provided.
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

	client := m.(*Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	// Authenticate with VMware Cloud Services
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	var sddc SDDC
	if sddcID != "" {
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		sddc, err = GetSddcByName(client, sddcName)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	// Check if already activated.
	if sddc.DeploymentStatus == constants.VmcActivationActiveStatus {
		return diag.Errorf("Already activated")
	}

	// Activate HCX.
	_, err = ActivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task to be completed.
	errcount := 0
	for {
		if sddcID != "" {
			sddc, err = GetSddcByID(client, sddcID)
		} else {
			sddc, err = GetSddcByName(client, sddcName)
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

		if sddc.DeploymentStatus == constants.VmcActivationActiveStatus {
			break
		}

		if sddc.DeploymentStatus == constants.VmcActivationFailedStatus {
			return diag.Errorf("Activation failed")
		}

		time.Sleep(constants.VmcRetryInterval)
	}

	return resourceVmcRead(ctx, d, m)
}

// resourceVmcRead retrieves the VMware Cloud on AWS configuration.
func resourceVmcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*Client)

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
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("****************")
	log.Printf("[Client inside]: %+v ", client)
	log.Printf("****************")

	var sddc SDDC
	if sddcID != "" {
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		sddc, err = GetSddcByName(client, sddcName)
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

	client := m.(*Client)

	token := client.Token
	sddcName := d.Get("sddc_name").(string)
	sddcID := d.Get("sddc_id").(string)

	// Authenticate with VMware Cloud Services
	accessToken, err := VmcAuthenticate(token)
	if err != nil {
		return diag.FromErr(err)
	}

	err = CloudAuthenticate(client, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	var sddc SDDC
	if sddcID != "" {
		sddc, err = GetSddcByID(client, sddcID)
	} else {
		sddc, err = GetSddcByName(client, sddcName)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// Deactivate HCX
	_, err = DeactivateHcxOnSDDC(client, sddc.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for task to be completed
	errcount := 0
	for {
		var sddc SDDC
		if sddcID != "" {
			sddc, err = GetSddcByID(client, sddcID)
		} else {
			sddc, err = GetSddcByName(client, sddcName)
		}
		if err != nil {
			// Attempt to bypass recurring situation where the HCX API
			// returns status 502 with a proxy server error, and an HTML response
			// instead of JSON.
			errcount++
			hclog.Default().Info("[INFO] - resourceVmcDelete() - Error retrieving SDDC status: ", "error", err.Error(), "Errcount:", errcount)
			if errcount > constants.VmcMaxRetries {
				return diag.FromErr(err)
			}
		}

		if sddc.DeploymentStatus == constants.VmcDeactivationInactiveStatus {
			break
		}

		if sddc.DeploymentStatus == constants.VmcDeactivationFailedStatus {
			return diag.Errorf("Deactivation failed")
		}

		time.Sleep(constants.VmcRetryInterval)
	}

	return diags
}
