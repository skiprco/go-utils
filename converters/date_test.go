package converters

import (
	"github.com/skiprco/go-utils/v2/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConvertDateForEn_success(t *testing.T) {
	result, genErr := ConvertToDate("2019-06-17T13:25:54.831Z", "EN")
	assert.Nil(t, genErr)
	assert.Equal(t, "17/06/2019", result)
}

func Test_ConvertDateForNL_success(t *testing.T) {
	result, genErr := ConvertToDate("2021-05-09T13:25:54.831Z", "NL")
	assert.Nil(t, genErr)
	assert.Equal(t, "09-05-2021", result)
}

func Test_ConvertDateForUnknownLanguage_success(t *testing.T) {
	result, genErr := ConvertToDate("2021-05-09T13:25:54.831Z", "KK")
	assert.Nil(t, genErr)
	assert.Equal(t, "09/05/2021", result)
}

func Test_ConvertDateForWithWrongFormat_success(t *testing.T) {
	result, genErr := ConvertToDate("2021-05-0913:25:54", "en")
	assert.Equal(t, errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorCannotConvertDate, map[string]string{"date": "2021-05-0913:25:54"}), genErr)
	assert.Equal(t, "", result)
}
