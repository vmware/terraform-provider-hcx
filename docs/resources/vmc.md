# Resource: `vmc`

This resource manages the activation and deactivation of HCX on VMC.

When HCX is activated, it is also configured with appropriate network and
compute profiles.

Ensure that the HCX appliances are reachable from the HCX connector for other
resources to work, (e.g. firewall configuration).

## Example Usage

```hcl
resource "hcx_vmc" "example" {
    sddc_name   = "example"
}

resource "hcx_site_pairing" "example" {
    url         = hcx_vmc.example.cloud_url
    username    = "cloudadmin@vmc.local"
    password    = var.vmc_vcenter_password
}
```

## Argument Reference

* `sddc_name` - (Optional) The name of the SDDC. Either `sddc_name` or `sddc_id`
  must be specified.
* `sddc_id` - (Optional) The ID of the SDDC. Either `sddc_id` or `sddc_name`
  must be specified.

## Attribute Reference

* `id` - The ID of the SDDC.
* `cloud_url` - The URL of HCX Cloud, used for the site pairing configuration.
* `cloud_type` - The type of the HCX Cloud. Should be `nsp` for VMC.
* `cloud_name` - The name of the HCX Cloud.
