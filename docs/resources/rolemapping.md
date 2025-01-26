# Resource: `rolemapping`

Assign the HCX roles to the user groups in vCenter that are allowed to perform
HCX operations.

## Example Usage

```hcl
resource "hcx_rolemapping" "rolemapping" {
    sso = hcx_sso.sso.id
    admin {
      user_group = "vsphere.local\\Administrators"
    }
    admin {
      user_group = "example.com\\Administrators"
    }
    enterprise {
      user_group = "example.com\\Administrators"
    }
}
```

## Argument Reference

* `sso` - (Required) The ID of the SSO Lookup Service.
* `admin` - (Optional) The group for `admin` users.
* `enterprise` - (Optional) The group for `enterprise` users.

### `admin` and `enterprise` Argument Reference

* `user_group` - (Optional) The group name. Defaults to `vsphere.local\\Administrators` for `admin`.
