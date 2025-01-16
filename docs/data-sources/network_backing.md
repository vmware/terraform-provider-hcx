# Data Source: `hcx_network_backing`

The `hcx_network_backing` data source retrieves information about a specific
network backing in an HCX environment.

## Example Usage

```hcl
data "hcx_network_backing" "example" {
  name         = "example"
  vcuuid       = "fa5482d8-a4f4-49a0-a4cd-3f6a99b2caee"
  network_type = "DistributedVirtualPortgroup"
}
```

## Argument Reference

* `name` - (Required) The name of the network backing.
* `vcuuid` - (Required) The UUID of the vCenter instance associated with the
  network backing.
* `network_type` - (Optional) The type of the network backing. Allowed values
  are `DistributedVirtualPortgroup` and `NsxtSegment`. Defaults to
  `DistributedVirtualPortgroup`.

## Attribute Reference

* `entityid` - The entity ID of the network backing.
