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

# This terraform DataSource is used to query the existing host from PowerMax array.
# The information fetched from this data source can be used for getting the details / for further processing in resource block.


# Returns all of the PowerMax hosts and their details
data "powermax_host" "HostDsAll" {
  # Optional Update the read timeout with (XXm) for minutes or (XXs) for timeout in seconds
  # If unset defaults to 2 minute timeout
  # timeouts = {
  #   read = "3m"
  # }
}

output "hostDsResultAll" {
  value = data.powermax_host.HostDsAll
}

# # Returns a subset of the PowerMax hosts based on the names provided in the `names` filter block and their details
# data "powermax_host" "HostDsFiltered" {
#   # Optional Update the read timeout with (XXm) for minutes or (XXs) for timeout in seconds
#   # If unset defaults to 2 minute timeout
#   # timeouts = {
#   #   read = "3m"
#   # }
#   filter {
#     # Optional list of names to filter upon
#     names = [
#       "Host124",
#       "Host173",
#     ]
#   }
# }

# output "hostDsResult" {
#   value = data.powermax_host.HostDsFiltered
# }

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
# Also, we can use the fetched information by the variable data.powermax_host.example
