# Data Source: `compute_profile`

The `compute_profile` data source retrieves the compute, storage, and network
settings that HCX uses on this site to deploy the Interconnect-dedicated virtual
appliances when a Service Mesh is added.

## Example Usage

```hcl
data "hcx_compute_profile" "vmc_cp" {
  vcenter = hcx_site_pairing.C2C1toC2C2.local_vc
  name    = "ComputeProfile(vcenter)"
}

output "compute_profile_vmc" {
  value = data.hcx_compute_profile.vmc_cp
}
```

## Argument Reference

* `name` - (Required) The name of the compute profile.
* `vcenter` - (Required) The ID of the vCenter instance.

## Attribute Reference

* `id` - ID of the compute profile.
