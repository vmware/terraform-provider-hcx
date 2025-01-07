# HCX Provider

This is the repository for the Terraform HCX Provider, which one can use with
Terraform to work with [VMware HCX][vmware-hcx].

[vmware-hcx]: https://cloud.vmware.com/vmware-hcx

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform


Documentation on HCX can be found at the [VMware HCX Documentation page](hhttps://docs.vmware.com/en/VMware-HCX/index.html)


# Using the Provider

The latest version of this provider requires Terraform v1.6.1 or higher to run.

Note that you need to run `terraform init` to fetch the provider before
deploying. Read about the provider split and other changes to TF v0.10.0 in the
official release announcement found [here][tf-0.10-announce].

[tf-0.10-announce]: https://www.hashicorp.com/blog/hashicorp-terraform-0-10/

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`provider` block in your Terraform configuration if you have not added one
already.

The syntax is as follows:

```hcl
provider "hcx" {
  version = "~> 1.0"
  ...
}
```


Version locking uses a pessimistic operator, so this version lock would mean
anything within the 1.x namespace, including or after 1.0.0. [Read
more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions

# Installation (automatic)

To install this provider, copy and paste this code into your Terraform configuration. Then, run terraform init.

```hcl
terraform {
  required_providers {
    hcx = {
      source = "vmware/hcx"
    }
  }
}

provider "hcx" {
  hcx       = "https://192.168.110.70"

  admin_username  = "admin"
  admin_password  = "changeme"

  username  = "administrator@corp.local"
  password  = "changeme"
}
```

# Installation (manual)


**NOTE:** Recommended way to compile the provider is using [Go Modules](https://blog.golang.org/using-go-modules), however vendored dependencies are still supported.

**NOTE:** For Terraform 0.13, please refer to [provider installation configuration][install-013] in order to use custom provider.
[install-013]: https://www.terraform.io/docs/commands/cli-config.html#provider-installation


## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/vmware/terraform-provider-hcx`:

```sh
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware
git clone https://github.com/vmware/terraform-provider-hcx.git
```

## Building and Installing the Provider

After the clone has been completed, you can enter the provider directory and build the provider.

```sh
cd $GOPATH/src/github.com/vmware/terraform-provider-hcx
make
```

After the build is complete, if your terraform running folder does not match your GOPATH environment, you need to copy the `terraform-provider-hcx` executable to your running folder and re-run `terraform init` to make terraform aware of your local provider executable.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Usage

In order to use the HCX Terraform provider you must first configure the provider to communicate with the HCX Connector. The HCX Connector is deployed at source (the cloud instance/the destination private cloud runs HCX Cloud).

Admin credentials are related to appliance configuration (appliance management on port 9443).
Normal credentials depends of the vCenter/SSO configuration made on the appliance.



## Example of Provider Configuration

```hcl
provider "hcx" {
  hcx       = "https://192.168.110.70"

  admin_username  = "admin"
  admin_password  = "changeme"

  username  = "administrator@corp.local"
  password  = "changeme"

  token     = "234567893245678345678" // Only needed if HCX on VMC on AWS SDDC need to be managed by TF
}
```

## Argument Reference

* `hcx` - (Optional) URL of the HCX connector. If not specified, only hcx_vmc is usable by this provider.
* `admin_username` - (Optional) Username of the HCX appliance. Only need if you want to manageable appliance setup.
* `admin_password` - (Optional) Password of the HCX appliance. Only need if you want to manageable appliance setup.
* `username` - (Optional) Username for HCX consumption. SSO/vSphere Role Mappings need to be set.
* `password` - (Optional) Password for HCX consumption. SSO/vSphere Role Mappings need to be set.
* `token` - (Required) VMware Cloud Service API Token. Generated from the VMware Cloud Services Console / My account / API Tokens. Environment variable VMC_API_TOKEN can be used to avoid setting the token in the code.


