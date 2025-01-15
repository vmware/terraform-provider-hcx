# Resource: `activation`

An activation key is mandatory to use a HCX system.

## Example Usage

```hcl
resource "hcx_activation" "activation" {
    activationkey = "*****-*****-*****-*****-*****"
}
```

## Argument Reference

* `activationkey` - (Required) The activation key.

## Attribute Reference

* `id` - The ID of the activation key.
