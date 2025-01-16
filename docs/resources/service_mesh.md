# Resource: `service_mesh`

An HCX Service Mesh is the effective HCX services configuration for a source and
destination site. A Service Mesh can be added to a connected Site Pair that has
a valid Compute Profile created on both sites.

Adding a Service Mesh initiates the deployment of HCX Interconnect virtual
appliances on both sites. An interconnect Service Mesh is always created at the
source site.

## Example Usage

```hcl
resource "hcx_service_mesh" "service_mesh_1" {
  name                          = "sm1"
  site_pairing                  = hcx_site_pairing.site1
  local_compute_profile         = hcx_compute_profile.compute_profile_1.name
  remote_compute_profile        = "Compute-RegionB01"
  app_path_resiliency_enabled   = false
  tcp_flow_conditioning_enabled = false
  uplink_max_bandwidth          = 10000
  service {
    name = "INTERCONNECT"
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
}

output "service_mesh_1" {
  value = hcx_service_mesh.service_mesh_1
}
```

## Argument Reference

* `name` - (Required) The name of the service mesh.
* `site_pairing` - (Required) The site pairing used by this service mesh.
* `local_compute_profile` - (Required) The local compute profile name.
* `remote_compute_profile` - (Required) The remote compute profile name.
* `app_path_resiliency_enabled` - (Optional) Enable the Application Path
  Resiliency feature. Defaults to `false`.
* `tcp_flow_conditioning_enabled` - (Optional) Enable the TCP flow conditioning
  feature. Defaults to `false`.
* `uplink_max_bandwidth` - (Optional) The maximum bandwidth used for uplinks.
  Defaults to `10000`.
* `service` - (Required) The list of HCX services. (Services selected here must
  be part of the compute profiles selected).
* `force_delete` - (Optional) Enable or disable force delete of the service
  mesh. Sometimes needed when site pairing is not connected anymore.
* `nb_appliances` - (Optional) The number of Network Extension appliances to
  deploy. Defaults to `1`.

### `service` Argument Reference

* `name` - (Required) The name of the HCX service. Allowed values include:
  `INTERCONNECT`, `WANOPT`, `VMOTION`, `BULK_MIGRATION`, `RAV`,
  `NETWORK_EXTENSION`, `DISASTER_RECOVERY`, or `SRM`.

## Attribute Reference

* `id` - ID of the Service Mesh.
