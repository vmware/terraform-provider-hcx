<!--
© Broadcom. All Rights Reserved.
The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
SPDX-License-Identifier: MPL-2.0
-->

<!-- markdownlint-disable first-line-h1 no-inline-html -->

<img src="docs/images/icon-color.svg" alt="VMware HCX" width="150">

# Terraform Provider for VMware HCX

[![Latest Release](https://img.shields.io/github/v/tag/vmware/terraform-provider-hcx?label=latest%20release&style=for-the-badge)](https://github.com/vmware/terraform-provider-hcx/releases/latest) [![License](https://img.shields.io/github/license/vmware/terraform-provider-hcx.svg?style=for-the-badge)](LICENSE)

The Terraform Provider for [VMware HCX][product-documentation] is a plugin for Terraform that allows you to interact with VMware HCX.

Learn more:

* Read the provider [documentation][provider-documentation].

* Join the community [discussions][provider-discussions].

## Requirements

* [Terraform 1.6+][terraform-install]

    For general information about Terraform, visit [HashiCorp Developer][terraform-install] and [the project][terraform-github] on GitHub.

* [Go 1.23.8][golang-install]

    Required, if [building][provider-build] the provider.

## Using the Provider

The Terraform Provider for VMware HCX is a Partner tier provider.

Partner tier providers are owned and maintained by a partner in the HashiCorp Technology Partner Program. HashiCorp verifies the authenticity of the publisher and the provider are listed on the [Terraform Registry][terraform-registry] with a Partner tier label.

To use a released version of the Terraform provider in your environment, run `terraform init` and Terraform will automatically install the provider from the Terraform Registry.

Unless you are contributing to the provider or require a pre-release bugfix or feature, use a
released version of the provider.

See [Installing the Terraform Provider for VMware HCX][provider-install] for additional instructions on automated and manual installation methods and how to control the provider version.

For either installation method, documentation about the provider configuration, resources, and data sources can be found on the Terraform Registry.

## Upgrading the Provider

The provider does not upgrade automatically. After each new release, you can run the following command to upgrade the provider:

```shell
terraform init -upgrade
```

## Contributing

The Terraform Provider for VMware HCX is the work of many contributors and the project team appreciates your help!

If you discover a bug or would like to suggest an enhancement, submit [an issue][provider-issues].

If you would like to submit a pull request, please read the [contribution guidelines][provider-contributing] to get started. In case of enhancement or feature contribution, we kindly ask you to open an issue to discuss it beforehand.

## Support

The Terraform Provider for VMware HCX is supported by Broadcom and the provider community. For bugs and feature requests please open a GitHub Issue and label it appropriately or contact Broadcom support.

## License

© Broadcom. All Rights Reserved.
The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.

The Terraform Provider for VMware HCX is available under the [Mozilla Public License, version 2.0][provider-license] license.

[golang-install]: https://golang.org/doc/install
[product-documentation]: https://docs.vmware.com/en/VMware-Cloud-on-AWS/index.html
[provider-contributing]: CONTRIBUTING.md
[provider-discussions]: https://github.com/vmware/terraform-provider-hcx/discussions
[provider-documentation]: https://registry.terraform.io/providers/vmware/hcx/latest/docs
[provider-build]: docs/build.md
[provider-install]: docs/install.md
[provider-issues]: https://github.com/vmware/terraform-provider-hcx/issues/new/choose
[provider-license]: LICENSE
[terraform-github]: https://github.com/hashicorp/terraform
[terraform-install]: https://developer.hashicorp.com/terraform/install
[terraform-registry]: https://registry.terraform.io

