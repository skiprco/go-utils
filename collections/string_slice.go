package collections

// StringSliceContains checks if a slice of strings contains the value
func StringSliceContains(slice []string, contains string) bool {
	for _, value := range slice {
		if value == contains {
			// Value found
			return true
		}
	}

	// Value not found
	return false
}
