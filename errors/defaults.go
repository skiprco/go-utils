package errors

// Defaults
var defaultMeta = map[string]string{}

// SetupDefaults sets defaults to created errors
func SetupDefaults(meta map[string]string) {
	// Set default meta
	if meta == nil {
		meta = map[string]string{}
	}
	defaultMeta = meta
}
