package metadata

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/errors"
)

// ========================================
// =                COMMON                =
// ========================================

const errDomain = "go-utils"
const errSubDomain = "metadata"

// Metadata is the generic representation of metadata.
// This type will be the same even when using different packages (gin, go-micro, ...)
type Metadata map[string]string

// ========================================
// =                  GOB                 =
// ========================================

// ToGob converts the metadata model to a Gob binary representation
func (meta Metadata) ToGob() []byte {
	// Ensure value is not nil
	if meta == nil {
		meta = Metadata{}
	}

	// Encode into Gob
	var metaBuffer bytes.Buffer
	enc := gob.NewEncoder(&metaBuffer)
	enc.Encode(meta) // Can only fail on nil value or nil pointer
	return metaBuffer.Bytes()
}

// FromGob creates a metadata model from its Gob binary representation
func FromGob(data []byte) (Metadata, *errors.GenericError) {
	// Return empty metadata object when no data is provided
	if data == nil || len(data) == 0 {
		return Metadata{}, nil
	}

	// Decode from Gob
	meta := &Metadata{}
	metaReader := bytes.NewReader(data)
	dec := gob.NewDecoder(metaReader)
	err := dec.Decode(meta)
	if err != nil {
		log.WithField("error", err).Error("Failed to decode Metadata from Gob bytes")
		return nil, errors.NewGenericError(400, errDomain, errSubDomain, "decode_metadata_from_glob_failed", nil)
	}
	return *meta, nil
}

// ========================================
// =                BASE64                =
// ========================================

// ToBase64 converts the metadata model to a base64 string
func (meta Metadata) ToBase64() string {
	// Ensure value is not nil
	if meta == nil {
		meta = Metadata{}
	}

	// Encode to base64
	return base64.StdEncoding.EncodeToString(meta.ToGob())
}

// FromBase64 creates a metadata model from its base64 string equivalent
func FromBase64(data string) (Metadata, *errors.GenericError) {
	// Return empty metadata object when no data is provided
	if data == "" {
		return Metadata{}, nil
	}

	// Decode from base64 to binary
	metaGob, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.WithField("error", err).WithField("data", data).Error("Failed to decode Gob bytes from base64 string")
		return nil, errors.NewGenericError(400, errDomain, errSubDomain, "decode_glob_from_base64_failed", nil)
	}

	// Decode from binary to Metadata
	return FromGob(metaGob)
}
