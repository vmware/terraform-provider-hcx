// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package constants

const (
	// HCX Cloud URLs
	HcxBaseURL          = "https://connect.hcx.vmware.com"
	HcxCloudURL         = HcxBaseURL
	HcxCloudAuthURL     = HcxBaseURL + "/provider/csp"
	HcxCloudConsumerURL = HcxBaseURL + "/provider/csp/consumer"

	// VMware Cloud URLs
	VmcBaseURL = "https://console.cloud.vmware.com"
	VmcAuthURL = VmcBaseURL + "/csp/gateway/am/api"

	// Network Types
	NetworkTypeDvpg       = "DistributedVirtualPortgroup"
	NetworkTypeNsxSegment = "NsxtSegment"
)

var AllowedNetworkTypes = []string{
	NetworkTypeDvpg,
	NetworkTypeNsxSegment,
}
