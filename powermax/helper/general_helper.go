// Copyright ©2023 Dell Inc. or its subsidiaries. All Rights Reserved.
package helper

// General Reusable Functions.
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
