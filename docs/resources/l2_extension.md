# Resource: `l2_extension`

You can bridge local network segments between HCX-enabled data centers with HCX
Network Extension.

With VMware HCX Network Extension (HCX-NE), a high-performance (4–6 Gbps)
service, you can extend the virtual machine networks to a VMware HCX-enabled
remote site. Virtual machines that are migrated or created on the extended
segment at the remote site are Layer 2 adjacent to virtual machines on the
origin network. Using Network Extension, a remote site's resources can be
quickly consumed. With Network Extension, the default gateway for the extended
network only exists at the source site. Traffic from virtual machines on remote
extended networks that must be routed returns to the source site gateway.

## Example Usage

```hcl
resource "hcx_l2_extension" "l2_extension_1" {
  site_pairing        = hcx_site_pairing.site1
  service_mesh_id     = hcx_service_mesh.service_mesh_1.id
  source_network      = "VM-RegionA01-vDS-COMP"
  NsxtSegment         = ""
  destination_t1      = "T1-GW"
  gateway             = "2.2.2.2"
  netmask             = "255.255.255.0"
  egress_optimization = false
  mon                 = true
  appliance_id        = hcx_service_mesh.service_mesh_1.appliances_id[1].id
}

output "l2_extension_1" {
  value = hcx_l2_extension.l2_extension_1
}
```

## Argument Reference

* `site_pairing` - (Required) The site pairing used by this service mesh.
* `service_mesh_id` - (Required) The ID of the Service Mesh to be used for this
  L2 extension.
* `source_network` - (Required) The source network. Must be a distributed port
  group which is VLAN tagged.
* `destination_t1` - (Required) The name of the NSX T1 at the destination.
* `gateway` - (Required) The gateway address to configure on the NSX T1. Should
  be equal to the existing default gateway at the source site.
* `netmask` - (Required) The netmask.
* `network_type` - (Optional) The network backing type. Allowed values include:
  `DistributedVirtualPortgroup` and `NsxtSegment`. Defaults to
  `DistributedVirtualPortgroup`.
* `appliance_id` - (Optional) The ID of the Network Extension appliance to use
  for the L2 extension. Defaults to the first appliance.
* `mon` - (Optional, default is false) Enable the MON (Mobility Optimized
  Networking) feature. Defaults to `false`.
* `egress_optimization` - (Optional, default is false) Enable the Egress
  Optimization feature. Defaults to `false`.

## Attribute Reference

* `id` - The ID of the L2 extension.
