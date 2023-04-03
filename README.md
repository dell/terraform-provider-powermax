# Terraform provider for PowerMax

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.1%20adopted-ff69b4.svg)](https://github.com/dell/terraform-provider-powermax/blob/main/about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/github/license/dell/terraform-provider-powermax)](https://github.com/dell/terraform-provider-powermax/blob/main/LICENSE)
[![Go version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://go.dev/dl/)
[![Terraform version](https://img.shields.io/badge/terraform-1.0+-blue.svg)](https://www.terraform.io/downloads)
[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/dell/terraform-provider-powermax?include_prereleases&label=latest&style=flat-square)](https://github.com/dell/terraform-provider-powermax/releases)


The Terraform Provider for PowerMax is a plugin for Terraform that allows the resource management of Powermax Storage arrays. For more details on PowerMax, please refer to PowerMax Official webpage [here][powermax-website].

For general information about Terraform, visit the [official website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io
[tf-github]: https://github.com/hashicorp/terraform
[powermax-website]: https://www.dell.com/en-in/dt/storage/powermax.htm?_gl=1*ji7vok*_ga*MTQ2NjY2MDI1Mi4xNjM0MTgzMzM3*_ga_1234567890*MTY2MDEwNzI4NC4xMC4wLjE2NjAxMDcyODQuMA..*_ga_5932KMEGPX*MTY2MDEwNzI4NC4xMC4wLjE2NjAxMDcyODQuNjA.&_ga=2.187158379.250612555.1660107285-1466660252.1634183337#tab0=0



## Table of Contents

* [Code of Conduct](about/CODE_OF_CONDUCT.md)
* [Committer Guide](about/COMMITTER_GUIDE.md)
* [Contributing Guide](about/CONTRIBUTING.md)
* [Developer](about/DEVELOPER.md)
* [Maintainers](about/MAINTAINERS.md)
* [Support](about/SUPPORT.md)
* [Attribution](about/ATTRIBUTION.md)
* [Additional Information](about/ADDITIONAL_INFORMATION.md)

## Supported Platforms
* PowerMax with Unisphere versions 10.0 and above.

## Prerequisites
* [Terraform >= 1.3.2](https://www.terraform.io)
* Go >= 1.19

## Installation
Install Terraform provider for PowerMax from terraform registry by adding the following block
```terraform
terraform {
  required_providers {
    powermax = { 
      version = "0.0.1"
      source = "registry.terraform.io/dell/powermax"
    }
  }
}
```


## Usage
Once you have installed the PowerMax provider, you can start using it in your Terraform configuration files. The provider has a number of resources that you can use to manage your PowerMax storage arrays.

For example, you can use the `powermax_storagegroup` resource to create a new storage group:
```terraform
resource "powermax_storagegroup" "test" {
  storage_group_id = "terraform_sg"
  srp_id           = "SRP_1"
  slo              = "Gold"
  host_io_limit = {
    host_io_limit_io_sec = "1000"
    host_io_limit_mb_sec = "1000"
    dynamicDistribution  = "Never"
  }
  volume_size            = "100"
  capacity_unit          = "CYL"
  volume_identifier_name = "terraform_volume"
  num_of_vols            = 1
}
```
For more resources, please refer to [examples](examples/main.tf)

## License
The Terraform provider for PowerMax is open-source software released under the [MPL-2.0 license](https://www.mozilla.org/en-US/MPL/2.0/).
