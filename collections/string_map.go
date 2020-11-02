package collections

// StringMapMerge creates a copy of the base map and merges the additional map.
// If both maps have the same key(s), the value of additional has priority.
// The returned map will never be nil, an empty map will be returned instead.
func StringMapMerge(base map[string]string, additional map[string]string) map[string]string {
	// Initialise result
	capacity := len(base) + len(additional)
	result := make(map[string]string, capacity)

	// Duplicate base
	for key, value := range base {
		result[key] = value
	}

	// Append additional
	for key, value := range additional {
		result[key] = value
	}

	// Return result
	return result
}
