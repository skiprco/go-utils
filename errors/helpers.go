package errors

// mergeMeta merges 2 meta maps
func mergeMeta(base map[string]string, additional map[string]string) map[string]string {
	// Return base if no additional is provided
	if additional == nil {
		return base
	}

	// Duplicate base
	result := map[string]string{}
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
