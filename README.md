<!--
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# Terraform Provider for Dell Technologies PowerMax
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-MPL_2.0-blue.svg)](LICENSE)

The Terraform Provider for Dell Technologies (Dell) PowerMax allows Data Center and IT administrators to use Hashicorp Terraform to automate and orchestrate the provisioning and management of Dell PowerMax storage systems.

The Terraform Provider can be used to manage volumes, storage groups, port groups, hosts, host groups and masking views.

## Table of Contents

* [Support](#support)
* [License](#license)
* [Prerequisites](#prerequisites)
* [List of DataSources in Terraform Provider for Dell PowerMax](#list-of-datasources-in-terraform-provider-for-dell-powermax)
* [List of Resources in Terraform Provider for Dell PowerMax](#list-of-resources-in-terraform-provider-for-dell-powermax)
* [Releasing, Maintenance and Deprecation](#releasing-maintenance-and-deprecation)
* [New to Terraform?](#new-to-terraform)

## Support
For any Terraform Provider for Dell PowerMax issues, questions or feedback, please follow our [support process](https://github.com/dell/dell-terraform-providers/blob/main/docs/SUPPORT.md)

## License
The Terraform Provider for Dell PowerMax is released and licensed under the MPL-2.0 license. See [LICENSE](LICENSE) for the full terms.

## Prerequisites

| **Terraform Provider** | **PowerMax Unisphere Version** | **OS**                                | **Terraform**    | **Golang** |
|------------------------|:-----------------------|:--------------------------------------|------------------|------------|
| v1.0.2                 | 10.0                           | ubuntu22.04 <br> rhel9.x <br> rhel8.x | 1.4.x <br> 1.5.x         | 1.20.x

## List of DataSources in Terraform Provider for Dell PowerMax
  * [Volume](docs/data-sources/volume.md)
  * [Storage Group](docs/data-sources/storagegroup.md)
  * [Port Group](docs/data-sources/portgroups.md)
  * [Host](docs/data-sources/host.md)
  * [Host Group](docs/data-sources/hostgroup.md)
  * [Masking View](docs/data-sources/maskingview.md)
  * [Port](docs/data-sources/port.md)
  * [Snapshot Policy](docs/data-sources/snapshotpolicy.md)
  * [Snapshot](docs/data-sources/snapshot.md)

## List of Resources in Terraform Provider for Dell PowerMax
  * [Volume](docs/resources/volume.md)
  * [Storage Group](docs/resources/storagegroup.md)
  * [Port Group](docs/resources/portgroup.md)
  * [Host](docs/resources/host.md)
  * [Host Group](docs/resources/hostgroup.md)
  * [Masking View](docs/resources/maskingview.md)
  * [Snapshot Policy](docs/resources/snapshotpolicy.md)
  * [Snapshot](docs/resources/snapshot.md)

## Installation and execution of Terraform Provider for Dell PowerMax
The installation and execution steps of Terraform Provider for Dell PowerMax can be found [here](about/INSTALLATION.md). 

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technnologies PowerMax follows [Semantic Versioning](https://semver.org/).

New versions will be release regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags in the form of "vx.y.z" where x.y.z corresponds to the version number.

## Documentation

For more detailed information, please refer to [Dell Terraform Providers Documentation](https://dell.github.io/terraform-docs/docs/storage/platforms/powermax/).

## New to Terraform?
**Here are some helpful links to get you started if you are new to terraform before using our provider:**

- Intro to Terraform: https://developer.hashicorp.com/terraform/intro 
- Providers: https://developer.hashicorp.com/terraform/language/providers 
- Resources: https://developer.hashicorp.com/terraform/language/resources
- Datasources: https://developer.hashicorp.com/terraform/language/data-sources
- Import: https://developer.hashicorp.com/terraform/language/import
- Variables: https://developer.hashicorp.com/terraform/language/values/variables
- Modules: https://developer.hashicorp.com/terraform/language/modules
- State: https://developer.hashicorp.com/terraform/language/state
