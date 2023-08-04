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
resource "powermax_storagegroup" "test" {
  name   = "terraform_sg"
  srp_id = "SRP_1"
  slo    = "Gold"
  host_io_limit = {
    host_io_limit_io_sec = "1000"
    host_io_limit_mb_sec = "1000"
    dynamic_distribution = "Never"
  }
  volume_ids = ["0008F"]
}