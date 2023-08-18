# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# v1.0.0
## Release Summary
The release supports PowerMax REST API 10.0.

## Changes made to json file
Removed volumeAttribute's required parameter num_of_vols Line 104797
Changed property name of storageGroupList to "name" from StorageGroupId Ln 98785, 98786, 98795. 98804
Change listVolumes filter parameters from type "array" to type "string" (Lines 52881 through 53308)

## Changes made to compile
.\model_rdf_group_label_list.go:26:15: undefined: RdfGroupID
.\model_rdf_group_label_list.go:79:47: undefined: RdfGroupID
.\model_rdf_group_label_list.go:81:13: undefined: RdfGroupID
.\model_rdf_group_label_list.go:89:50: undefined: RdfGroupID
.\model_rdf_group_label_list.go:106:47: undefined: RdfGroupID
Changed to RdfGroupId
.\model_volume.go:839:10: invalid operation: cannot indirect o.HasEffectiveWwn (value of type func() bool)
.\model_volume.go:848:9: cannot use o.HasEffectiveWwn (value of type func() bool) as type *bool in return statement
.\model_volume.go:862:2: cannot assign to o.HasEffectiveWwn (value of type func() bool)
.\model_volume.go:884:18: field and method with the same name HasEffectiveWwn
       .\model_volume.go:70:2: other declaration of HasEffectiveWwn
Line 884 - comment the function HasEffectiveWwn

Update client.go to comment out unused APIs

