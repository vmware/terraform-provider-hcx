<img src="https://raw.githubusercontent.com/vmware/terraform-provider-hcx/main/docs/images/icon-color.svg" alt="VMware HCX" width="150">

# Terraform Provider for VMware HCX

The Terraform Provider for [VMware HCX][product-documentation] is a plugin for Terraform that allows you to
interact with VMware HCX.

## Example Usage

```hcl
provider "hcx" {
  hcx            = "https://sfo-hcx01.example.com"
  admin_username = "admin"
  admin_password = "VMware1!"
  username       = "svc-hcx@example.com"
  password       = "VMware1!"
  token          = "123456789123456789" // Only needed for HCX on VMware Cloud on AWS.
}
```

## Argument Reference

* `hcx` - (Optional) URL of the HCX connector. If not specified, only `hcx_vmc` is usable by this provider.
* `admin_username` - (Optional) Username of the HCX appliance. Only need if you want to manage the appliance setup.
* `admin_password` - (Optional) Password of the HCX appliance. Only need if you want to manage the appliance setup.
* `username` - (Optional) Username for HCX consumption. SSO/vSphere Role Mappings need to be set.
* `password` - (Optional) Password for HCX consumption. SSO/vSphere Role Mappings need to be set.
* `token` - (Required) VMware Cloud Service API Token. Generated from the **VMware Cloud Services Console** > **My account** > **API Tokens**. Environment variable `VMC_API_TOKEN` can be used to avoid setting the token in the code.

[product-documentation]: https://techdocs.broadcom.com/us/en/vmware-cis/hcx.html
