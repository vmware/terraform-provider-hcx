# Resource: `compute_profile`

A Compute Profile contains the compute, storage, and network settings that HCX
uses on this site to deploy the Interconnect-dedicated virtual appliances when a
Service Mesh is added.

Create a Compute Profile in the Multi-Site Service Mesh interface in both the
source and the destination HCX environments using the planned configuration
options for each site, respectively.

## Example Usage

```hcl
resource "hcx_compute_profile" "compute_profile_1" {
  name                = "comp1"
  datacenter          = "RegionA01-ATL"
  cluster             = "RegionA01-COMP01"
  datastore           = "RegionA01-ISCSI01-COMP01"
  management_network  = hcx_network_profile.net_management.id
  replication_network = hcx_network_profile.net_management.id
  uplink_network      = hcx_network_profile.net_uplink.id
  vmotion_network     = hcx_network_profile.net_vmotion.id
  dvs                 = "RegionA01-vDS-COMP"
  service {
    name = "INTERCONNECT"
  }
  service {
    name = "WANOPT"
  }
  service {
    name = "VMOTION"
  }
  service {
    name = "BULK_MIGRATION"
  }
  service {
    name = "RAV"
  }
  service {
    name = "NETWORK_EXTENSION"
  }
  service {
    name = "DISASTER_RECOVERY"
  }
  service {
    name = "SRM"
  }
}

output "compute_profile_1" {
  value = hcx_compute_profile.compute_profile_1
}
```

## Argument Reference

* `name` - (Required) The name of the compute profile.
* `datacenter` - (Required) The datacenter where HCX services will be available.
* `cluster` - (Required) The cluster used for HCX appliances deployment.
* `datastore` - (Required) The datastore used for HCX appliances deployment.
* `management_network` - (Required) The management network profile (ID).
* `replication_network` - (Required) The replication network profile (ID).
* `vmotion_network` - (Required) The vMotion network profile (ID).
* `uplink_network` - (Required) The uplink network profile (ID).
* `dvs` - (Required) The distributed switch used for L2 extension.
* `service` - (Required) The list of HCX services.

### `service` Argument Reference

* `name` - (Required) The name of the HCX service. Allowed values include:
  `INTERCONNECT`, `WANOPT`, `VMOTION`, `BULK_MIGRATION`, `RAV`,
  `NETWORK_EXTENSION`, `DISASTER_RECOVERY`, or `SRM`.

## Attribute Reference

* `id` - ID of the compute profile.
