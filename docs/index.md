<img src="https://raw.githubusercontent.com/vmware/terraform-provider-hcx/main/docs/images/icon-color.svg" alt="VMware HCX" width="150">

# Terraform Provider for VMware HCX

The Terraform Provider for [VMware HCX][product-documentation] is a plugin for
Terraform that allows you to interact with VMware HCX.

## Example Usage

```hcl
provider "hcx" {
  hcx            = "https://sfo-hcx01.example.com"
  admin_username = "admin"
  admin_password = "VMware1!"
  username       = "svc-hcx@example.com"
  password       = "VMware1!"
  vmc_token      = "123456789123456789" // Only needed for HCX on VMware Cloud on AWS.
}
```

## Argument Reference

* `hcx` - (Optional) URL of the HCX connector. If not specified, only `hcx_vmc`
  is usable by this provider. Can also be specified with the `HCX_URL`
  environment variable.
* `admin_username` - (Optional) Username of the HCX appliance. Can also be
  specified with the `HCX_USER` environment variable.
* `admin_password` - (Optional) Password of the HCX appliance. Can also be
  specified with the `HCX_PASSWORD` environment variable.
* `username` - (Optional) Username for HCX consumption. Can also be specified
  with the `HCX_ADMIN_USER` environment variable.
* `password` - (Optional) Password for HCX consumption. Can also be specified
  with the `HCX_ADMIN_PASSWORD` environment variable.
* `allow_unverified_ssl` - (Optional) Allow SSL connections with unverifiable
  certificates. Defaults to `false`. Can also be specified with the
  `HCX_ALLOW_UNVERIFIED_SSL` environment variable.
* `vmc_token` - (Optional) VMware Cloud Service API Token. This token is
  generated from the **VMware Cloud Services Console** > **My Account** > **API
  Tokens**. Can also be specified with the `VMC_API_TOKEN` environment variable.

[product-documentation]: https://techdocs.broadcom.com/us/en/vmware-cis/hcx.html
