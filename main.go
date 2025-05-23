// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/vmware/terraform-provider-hcx/hcx"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// main initializes and starts the plugin service for the provider with optional debugging support.
func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return hcx.Provider()
		},
	}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "vmware/hcx"
	}

	plugin.Serve(opts)
}
