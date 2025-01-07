---
page_title: "HCX/VMC Lab - Full HCX Connector configuration"
---

This TF file automates the configuration of a HCX lab. It manages creation/update/deletion of :
* HCX activation and configuration at VMC side
* Site pairing
* Network profiles (Management, vMotion and Uplink)
* Compute profile
* Service Mesh
* L2 Extension

## Usage example

```hcl

terraform {
  required_providers {
    hcx = {
      source = "vmware/hcx"
    }
  }
}

provider hcx {
    hcx             = "https://172.17.9.10"

    admin_username  = "admin"
    admin_password  = "VMware1!VMware1!"

    username        = "administrator@vsphere.local"
    password        = "VMware1!"

    // token = "xxx" if VMC
    // export VMC_API_TOKEN=...
}

// Variables definitions
variable "hcx_vmc_vcenter_password" {
  type        = string
  description = "vCenter password (export TF_VAR_hcx_vmc_vcenter_password=...)"
}

// Provider config
provider hcx {
    hcx             = "https://172.17.9.10"

    admin_username  = "admin"
    admin_password  = "VMware1!VMware1!"

    username        = "administrator@cpod-vcn.az-fkd.cloud-garage.net"
    password        = "changeme"
}

resource "hcx_vcenter" "vcenter" {
    url         = "https://172.17.9.3"
    username    = "administrator@cpod-vcn.az-fkd.cloud-garage.net"
    password    = "changeme"

    depends_on  = [hcx_activation.activation]
}

resource "hcx_sso" "sso" {
    vcenter     = hcx_vcenter.vcenter.id
    url         = "https://172.17.9.3"
}


resource "hcx_rolemapping" "rolemapping" {
    sso = hcx_sso.sso.id

    admin {
      user_group = "cpod-vcn.az-fkd.cloud-garage.net\\Administrators"
    }

    admin {
      user_group = "cpod-vcn.az-fkd.cloud-garage.net\\Administrators"
    }

    enterprise {
      user_group = "cpod-vcn.az-fkd.cloud-garage.net\\Administrators"
    }
}

resource "hcx_location" "location" {
    city        = "Paris"
    country     = "France"
    province    = "Ile-de-France"
    latitude    = 48.86669293
    longitude   = 2.333335326
}



// Datasources and Resources
resource "hcx_vmc" "vmc_nico" {  
    sddc_name   = "mySDDC-name"
}

resource "hcx_site_pairing" "vmc" {
    url         = hcx_vmc.vmc_nico.cloud_url
    username    = "cloudadmin@vmc.local"
    password    = var.vmc_vcenter_password

    depends_on  = [hcx_rolemapping.rolemapping]
}

resource "hcx_network_profile" "net_management" {
  vcenter       = hcx_site_pairing.vmc.local_vc
  network_name  = "HCX-Management-RegionA01"
  name          = "HCX-Management-RegionA01-profile"
  mtu           = 1500

  ip_range {
    start_address   = "192.168.110.151"
    end_address     = "192.168.110.155"
  }

  gateway           = "192.168.110.1"
  prefix_length     = 24
  primary_dns       = "192.168.110.10"
  secondary_dns     = ""
  dns_suffix        = "corp.local"
}



resource "hcx_network_profile" "net_uplink" {
  vcenter       = hcx_site_pairing.vmc.local_vc
  network_name  = "HCX-Uplink-RegionA01"
  name          = "HCX-Uplink-RegionA01-profile"
  mtu           = 1600

  ip_range {
    start_address   = "192.168.110.156"
    end_address     = "192.168.110.160"
  }

  gateway           = "192.168.110.1"
  prefix_length     = 24
  primary_dns       = "192.168.110.1"
  secondary_dns     = ""
  dns_suffix        = "corp.local"
}

resource "hcx_network_profile" "net_vmotion" {
  vcenter       = hcx_site_pairing.vmc.local_vc
  network_name  = "HCX-vMotion-RegionA01"
  name          = "HCX-vMotion-RegionA01-profile"
  mtu           = 1500

  ip_range {
    start_address   = "10.10.30.151"
    end_address     = "10.10.30.155"
  }

  gateway           = ""
  prefix_length     = 24
  primary_dns       = ""
  secondary_dns     = ""
  dns_suffix        = ""
}



resource "hcx_compute_profile" "compute_profile_1" {
  name                  = "comp1"
  datacenter            = "RegionA01-ATL"
  cluster               = "RegionA01-COMP01"
  datastore             = "RegionA01-ISCSI01-COMP01"

  management_network    = hcx_network_profile.net_management.id
  replication_network   = hcx_network_profile.net_management.id
  uplink_network        = hcx_network_profile.net_uplink.id
  vmotion_network       = hcx_network_profile.net_vmotion.id
  dvs                   = "RegionA01-vDS-COMP"

  service {
    name                = "INTERCONNECT"
  }

  service {
    name                = "WANOPT"
  }

  service {
    name                = "VMOTION"
  }

  service {
    name                = "BULK_MIGRATION"
  }

  service {
    name                = "NETWORK_EXTENSION"
  }

  service {
    name                = "DISASTER_RECOVERY"
  }

}

resource "hcx_service_mesh" "service_mesh_1" {
  name                            = "sm1"
  site_pairing                    = hcx_site_pairing.vmc
  local_compute_profile           = hcx_compute_profile.compute_profile_1.name
  remote_compute_profile          = "ComputeProfile(vcenter)"

  app_path_resiliency_enabled     = false
  tcp_flow_conditioning_enabled   = false

  uplink_max_bandwidth            = 10000

  service {
    name                = "INTERCONNECT"
  }

  service {
    name                = "VMOTION"
  }

  service {
    name                = "BULK_MIGRATION"
  }

  service {
    name                = "NETWORK_EXTENSION"
  }

  service {
    name                = "DISASTER_RECOVERY"
  }

}


```
