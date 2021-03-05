package metadata

const errorDomain = "go_utils"
const errorSubDomain = "metadata"

// ErrorUserIDNotInMeta indicates we tried to extract the user
// from the metadata, but "user_id" is not set.
const ErrorUserIDNotInMeta = "user_id_not_set_in_metadata"

// ErrorDecodeGlobFromBase64Failed indicates decoding the metadata as glob
// from the provided base64 string failed
const ErrorDecodeGlobFromBase64Failed = "decode_glob_from_base64_failed"
