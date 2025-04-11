# Resource: `vcenter`

The vCenter instance used by the HCX system. If the resource is created or
updated, the application service is restarted.

## Example Usage

```hcl
resource "hcx_vcenter" "vcenter" {
  url        = "https://vcsa-01a.example.com"
  username   = "administrator@vsphere.local"
  password   = "VMware1!"
  depends_on = [hcx_activation.activation]
}
```

## Argument Reference

* `url` - (Required) The URL of the vCenter instance.
* `username` - (Required) The username to authenticate to the vCenter instance.
* `password` - (Required) The password to authenticate to the vCenter instance.

## Attribute Reference

* `id` - The UUID of the vCenter instance.
