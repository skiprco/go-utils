package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringSliceContains_SliceNil(t *testing.T) {
	result := StringSliceContains(nil, "")
	assert.False(t, result)
}

func Test_StringSliceContains_Found(t *testing.T) {
	slice := []string{"test1", "test2"}
	result := StringSliceContains(slice, "test1")
	assert.True(t, result)
}

func Test_StringSliceContains_NotFound(t *testing.T) {
	slice := []string{"test1", "test2"}
	result := StringSliceContains(slice, "test3")
	assert.False(t, result)
}
