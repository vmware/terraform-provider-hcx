# Resource: `site_pairing`

A Site Pair establishes the connection needed for management, authentication,
and orchestration of HCX services across a source and destination environment.

In HCX Connector to HCX Cloud deployments, the HCX Connector is deployed at the
legacy or source vSphere environment. The HCX Connector creates a unidirectional
site pairing to an HCX Cloud system. In this type of site pairing, all HCX
Service Mesh connections, migration, and network extension operations, including
reverse migrations, are always initiated from the HCX Connector at the source.

## Example Usage

```hcl
resource "hcx_site_pairing" "site1" {
  url      = "https://hcx-cloud-01b.example.com"
  username = "administrator@vsphere.local"
  password = "VMware1!"
}

output "hcx_site_pairing_site1" {
  value = hcx_site_pairing.site1
}
```

## Argument Reference

* `url` - (Required) The URL of the remote cloud.
* `username` - (Required) The username used for remote cloud authentication.
* `password` - (Required) The password used for remote cloud authentication.

## Attribute Reference

* `id` - The ID of the site pairing.
* `local_vc` - The ID of the local vCenter instance.
* `local_endpoint_id` - The endpoint ID of the local HCX site.
* `local_name` - The endpoint name of the local HCX site.
* `remote_name` - The endpoint name of the remote HCX site.
* `remote_endpoint_type` - The endpoint type of the remote HCX site.
* `remote_resource_id` - The resource ID of the remote HCX site.
* `remote_resource_name` - The resource name of the remote HCX site.
* `remote_resource_type` - The resource type of the remote HCX site.
