package core

import "slices"

// MapKeysSorted returns unsorted map keys as a slice or nil for nil maps
func MapKeys[TValue any](values map[string]TValue) []string {
	if values == nil {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

// MapKeysSorted returns alphabetically sorted map keys as a slice or nil for nil maps
func MapKeysSorted[TValue any](values map[string]TValue) []string {
	if values == nil {
		return nil
	}
	sortedKeys := MapKeys(values)
	slices.Sort(sortedKeys)
	return sortedKeys
}
