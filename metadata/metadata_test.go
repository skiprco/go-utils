package metadata

import (
	"testing"

	"github.com/skiprco/go-utils/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// =              ROUNDTRIPS              =
// ========================================

func Test_GobRoundtrip(t *testing.T) {
	meta := fixtureMetadata()
	result, genErr := FromGob(meta.ToGob())
	require.Nil(t, genErr)
	assert.Equal(t, meta, result)
}

func Test_Base64Roundtrip(t *testing.T) {
	meta := fixtureMetadata()
	result, genErr := FromBase64(meta.ToBase64())
	require.Nil(t, genErr)
	assert.Equal(t, meta, result)
}

// ========================================
// =                  GOB                 =
// ========================================
func Test_ToGob_Nil(t *testing.T) {
	// Should default to empty Metadata
	var meta Metadata = nil
	result, genErr := FromGob(meta.ToGob())
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, result)
}

func Test_FromGob_DataNil(t *testing.T) {
	// Should default to empty Metadata
	result, genErr := FromGob(nil)
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, result)
}

func Test_FromGob_DataEmpty(t *testing.T) {
	// Should default to empty Metadata
	result, genErr := FromGob([]byte{})
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, result)
}

func Test_FromGob_DataInvalid(t *testing.T) {
	// Should default to empty Metadata
	result, genErr := FromGob([]byte("invalid"))
	assert.Nil(t, result)
	errors.AssertGenericError(t, genErr, 400, "decode_metadata_from_glob_failed", nil)
}

// ========================================
// =                BASE64                =
// ========================================
func Test_ToBase64_Nil(t *testing.T) {
	// Should default to empty Metadata
	var meta Metadata = nil
	result, genErr := FromBase64(meta.ToBase64())
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, result)
}

func Test_FromBase64_DataEmpty(t *testing.T) {
	// Should default to empty Metadata
	result, genErr := FromBase64("")
	require.Nil(t, genErr)
	assert.Equal(t, Metadata{}, result)
}

func Test_FromBase64_DataInvalid(t *testing.T) {
	// Should default to empty Metadata
	result, genErr := FromBase64("invalid")
	assert.Nil(t, result)
	errors.AssertGenericError(t, genErr, 400, "decode_glob_from_base64_failed", nil)
}
