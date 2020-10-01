package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeMeta_AdditionalNil(t *testing.T) {
	// Setup test data
	base := map[string]string{
		"test1": "success1",
		"test2": "failed2",
	}

	// Call function
	result := mergeMeta(base, nil)

	// Assert result
	assert.NotSame(t, result, base)
	expected := map[string]string{
		"test1": "success1",
		"test2": "failed2",
	}
	assert.Equal(t, expected, result)
}

func Test_mergeMeta_AdditionalSet(t *testing.T) {
	// Setup test data
	base := map[string]string{
		"test1": "success1",
		"test2": "failed2",
	}

	addition := map[string]string{
		"test2": "success2",
		"test3": "success3",
	}

	// Call function
	result := mergeMeta(base, addition)

	// Assert result
	assert.NotSame(t, result, base)
	assert.NotSame(t, result, addition)
	expected := map[string]string{
		"test1": "success1",
		"test2": "success2",
		"test3": "success3",
	}
	assert.Equal(t, expected, result)
}
