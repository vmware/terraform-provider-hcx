# Resource: `network_profile`

The Network Profile is an abstraction of a distributed port group, standard port
group, or NSX logical switch, and the Layer 3 properties of that network. A
Network Profile is a subcomponent of a complete Compute Profile.

Create a Network Profile for each network you intend to use with the HCX
services. The extension selects these network profiles when creating a Compute
Profile and assigns one or more of four Network Profile functions.

## Example Usage

```hcl
resource "hcx_network_profile" "net_management" {
  vcenter      = hcx_site_pairing.site1.local_vc
  network_name = "HCX-Management-RegionA01"
  name         = "HCX-Management-RegionA01-profile"
  mtu          = 1500
  ip_range {
    start_address = "192.168.110.151"
    end_address   = "192.168.110.155"
  }
  gateway       = "192.168.110.1"
  prefix_length = 24
  primary_dns   = "192.168.110.10"
  secondary_dns = ""
  dns_suffix    = "example.com"
}

output "net_management" {
  value = hcx_network_profile.net_management
}
```

## Example Usage (VMC)

```hcl
resource "hcx_network_profile" "net_management" {
  vcenter = hcx_site_pairing.C2C1toC2C2.local_vc
  vmc     = true
  name    = "externalNetwork"
  mtu     = 1500
  ip_range {
    start_address = "18.132.147.242"
    end_address   = "18.132.147.242"
  }
  ip_range {
    start_address = "18.168.66.74"
    end_address   = "18.168.66.74"
  }
  prefix_length = 0
}
```

## Argument Reference

* `site_pairing` - (Required) The site pairing map, to be retrieved with the
  `hcx_site_pairing` resource.
* `network_name` - (Required) The network name used for the profile.
* `name` - (Required) The name of the network profile.
* `mtu` - (Required) The MTU of the network profile.
* `gateway` - (Optional) The gateway for the network profile.
* `prefix_length` - (Required) The prefix length for the network profile.
* `primary_dns` - (Optional) The primary DNS for the network profile.
* `secondary_dns` - (Optional) The secondary DNS for the network profile.
* `dns_suffix` - (Optional) The DNS suffix for the network profile.
* `ip_range` - (Required) The list of IP ranges.
* `vmc` - (Optional) If set to true, the network profile will not be created or
  deleted, only IP pools will be updated.

### `ip_range` Argument Reference

* `start_address` - (Required) The start address of the IP pool for the network
  profile.
* `end_address` - (Required) The end address of the IP pool for the network
  profile.

## Attribute Reference

* `id` - The ID of the network profile.
