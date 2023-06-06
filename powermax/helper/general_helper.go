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

package helper

// IsParamUpdated General Reusable Functions.
func IsParamUpdated(updatedParams []string, paramName string) bool {
	isParamUpdate := false
	for _, updatedParam := range updatedParams {
		if updatedParam == paramName {
			isParamUpdate = true
			break
		}
	}
	return isParamUpdate
}

// CompareStringSlice Compare string slices. return true if the length and elements are same.
func CompareStringSlice(planInitiators, stateInitiators []string) bool {
	if len(planInitiators) != len(stateInitiators) {
		return false
	}

	itemAppearsTimes := make(map[string]int, len(planInitiators))
	for _, i := range planInitiators {
		itemAppearsTimes[i]++
	}

	for _, i := range stateInitiators {
		if _, ok := itemAppearsTimes[i]; !ok {
			return false
		}

		itemAppearsTimes[i]--
		if itemAppearsTimes[i] == 0 {
			delete(itemAppearsTimes, i)
		}
	}
	return len(itemAppearsTimes) == 0
}
