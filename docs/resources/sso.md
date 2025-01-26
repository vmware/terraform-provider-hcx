# Resource: `sso`

The vCenter SSO instance used by the HCX system. If the resource is created or
updated, the application service is restarted.

## Example Usage

```hcl
resource "hcx_vcenter" "vcenter" {
  url        = "https://vc01.example.com"
  username   = "administrator@vsphere.local"
  password   = "VMware1!"
  depends_on = [hcx_activation.activation]
}

resource "hcx_sso" "sso" {
  vcenter = hcx_vcenter.vcenter.id
  url     = hcx_vcenter.vcenter.url
}
```

## Argument Reference

* `vcenter` - (Required) The ID of the vCenter instance.
* `url` - (Required) The URL of the vCenter instance.

## Attribute Reference

* `id` - The UUID of the vCenter instance.
