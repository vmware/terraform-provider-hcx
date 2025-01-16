# Resource: `location`

Select the nearest major city to where the HCX system is geographically located.
HCX sites are represented visually in the Dashboard.

Select the nearest major city to where the HCX system is geographically located.
HCX sites are visually represented in the Dashboard.

## Example Usage

```hcl
resource "hcx_location" "location" {
    city        = "Paris"
    country     = "France"
    province    = "Ile-de-France"
    latitude    = 48.86669293
    longitude   = 2.333335326
}
```

## Argument Reference

* `city` - (Optional) The city where the HCX site is located.
* `province` - (Optional) The province where the HCX site is located.
* `country` - (Optional) The country where the HCX site is located.
* `latitude` - (Optional) The latitude coordinate of the HCX site.
* `longitude` - (Optional) The longitude coordinate of the HCX site.
