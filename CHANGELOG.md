# CHANGELOG

## [v0.5.1](https://github.com/vmware/terraform-provider-hcx/releases/tag/v0.5.1)

> Release Date: 2025-01-15

General:

- Transferred to [`vmware/terraform-provider-hcx`](https://github.com/vmware/terraform-provider-hcx) from [`adeleporte/terraform-provider-hcx`](https://github.com/adeleporte/terraform-provider-hcx).
- License changed from Apache License 2.0 to Mozilla Public License 2.0. [#17](https://github.com/vmware/terraform-provider-hcx/pull/17)

Chores:

- Updated Go to v1.22.7. [#47](https://github.com/vmware/terraform-provider-hcx/pull/47)
- Removes the use of `io/ioutil`. As of Go 1.16, the same functionality is now provided by package `io` or package `os`, and those implementations are preferred. [#48](https://github.com/vmware/terraform-provider-hcx/pull/48)
