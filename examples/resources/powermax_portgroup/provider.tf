/*
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
*/
terraform {
  required_providers {
    powermax = {
      source = "dell/powermax"
    }
  }
}

provider "powermax" {
  username      = var.username
  password      = var.password
  endpoint      = var.endpoint
  serial_number = var.serial_number
  pmax_version  = var.pmax_version
  insecure      = true

  ## Provider can also be set using environment variables
  ## If environment variables are set it will override this configuration
  ## Example environment variables
  # POWERMAX_USERNAME="username"
  # POWERMAX_PASSWORD="password"
  # POWERMAX_ENDPOINT="https://yourhost.host.com:8443"
  # POWERMAX_SERIAL_NUMBER="xxxxxxxxxxxx"
  # POWERMAX_POWERMAX_VERSION="100"
  # POWERMAX_INSECURE="false"
}