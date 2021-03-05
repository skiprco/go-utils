package manifest

const errorDomain = "go_utils"
const errorSubDomain = "manifest"

// ErrorManifestFileNotFound indicates the manifest file is not found or is not readable.
const ErrorManifestFileNotFound = "manifest_file_not_found"

// ErrorUnmarshalManifestFailed indicates parsing the manifest file as JSON failed.
const ErrorUnmarshalManifestFailed = "unmarshal_manifest_failed"
