// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package constants

import "time"

const (

	// HCX
	HcxDefaultPort = 9443

	// APIs
	HcxActivationAPI         = "/api/admin/global/config/hcx"
	HcxAppEngineBaseAPI      = "/components/appengine"
	HcxAppEngineStartAPI     = HcxAppEngineBaseAPI + "?action=start"
	HcxAppEngineStopAPI      = HcxAppEngineBaseAPI + "?action=stop"
	HcxAppEngineStatusAPI    = HcxAppEngineBaseAPI + "/status"
	HcxCertificatesAPI       = "/hybridity/api/admin/certificates"
	HcxComputeProfileAPI     = "/hybridity/api/interconnect/computeProfiles"
	HcxL2ExtensionsBaseAPI   = "/hybridity/api/l2Extensions"
	HcxL2ExtensionsGetAPI    = HcxL2ExtensionsBaseAPI
	HcxL2ExtensionsCreateAPI = HcxL2ExtensionsBaseAPI
	HcxL2ExtensionsDeleteAPI = HcxL2ExtensionsBaseAPI + "/%s"

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

	// VMC
	VmcMaxRetries                 = 12
	VmcRetryInterval              = 10 * time.Second
	VmcActivationActiveStatus     = "ACTIVE"
	VmcActivationFailedStatus     = "ACTIVATION_FAILED"
	VmcDeactivationInactiveStatus = "DE-ACTIVATED"
	VmcDeactivationFailedStatus   = "DEACTIVATION_FAILED"

	// Status
	StoppedStatus  = "STOPPED"
	RunningStatus  = "RUNNING"
	FailedStatus   = "FAILED"
	SuccessStatus  = "SUCCESS"
	RealizedStatus = "REALIZED"

	// Network Profile
	DefaultNetworkProfileOrg = "DEFAULT"

	// Location
	DefaultLatitude  = 0
	DefaultLongitude = 0

	// Compute Profile
	DefaultComputeType = "VC"

	// Single Sign-On
	DefaultSsoProviderType = "PSC"

	// Role Mappings
	RoleSystemAdmin     = "System Administrator"
	RoleEnterpriseAdmin = "Enterprise Administrator"
)

var AllowedNetworkTypes = []string{
	NetworkTypeDvpg,
	NetworkTypeNsxSegment,
}
