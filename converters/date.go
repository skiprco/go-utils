package converters

import (
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/skiprco/go-utils/v2/errors"
	"time"
)

// based on https://en.wikipedia.org/wiki/Date_format_by_country
// sadly there is currently no go library to manage that. Probably not the best way to do that, but it should works until
// Skipr a transcontinental company
var dateFormats = map[string]dateFormat{
	"FR":      {date: "02/01/06"},
	"NL":      {date: "02-01-06"},
	"EN":      {date: "02/01/06"},
	"DEFAULT": {date: "02/01/06"},
}

type dateFormat struct {
	date string
}

// ConvertToDate converts a string with a date of ISO format to a human readable date string formatted by the expected language
//
// Raises
//
// not error can be returned
func ConvertToDate(date, lang string) (string, *errors.GenericError) {
	format, found := dateFormats[lang]
	if !found {
		format = dateFormats["DEFAULT"]
	}
	printer := getPrinter(lang)
	t, err := time.Parse(protocol.ISO8601TimeFormat, date)
	if err != nil {
		meta := map[string]string{
			"date": date,
		}
		return "", errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorCannotConvertDate, meta)
	}
	return printer.Sprint(t.Format(format.date)), nil
}
