// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"fmt"
	"github.com/vmware/terraform-provider-hcx/hcx/constants"
)

// ValidateNetworkType validates that the provided value is a string and matches one of the allowed network types.
// Returns warnings and errors based on value validation.
func ValidateNetworkType(val interface{}, key string) (warns []string, errs []error) {
	networkType, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string, got: %T", key, val))
		return warns, errs
	}

	for _, allowedType := range constants.AllowedNetworkTypes {
		if networkType == allowedType {
			return warns, errs
		}
	}

	errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, constants.AllowedNetworkTypes, networkType))
	return warns, errs
}
