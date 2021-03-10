package converters

import (
	"github.com/skiprco/go-utils/v2/errors"
)

// ConvertCentToCurrency converts a cent value into a human readable price with the expected currency
// decimal are managed by the expected lang
// lang should have the ISO 639-1 format ("fr", "en"). If the language is not found or empty, use English
// the currency should be the ISO 4217 format ('EUR','USD'). If the currency doesn't implemented yet, return an error unknown_currency
//
// Raises
//
// - 500/unknown_currency: The currency is unknown. If it's happened, extend the switch with the expected currency
func ConvertCentToCurrency(cent int64, currency, lang string) (string, *errors.GenericError) {
	printer := getPrinter(lang)
	value := convertCentToUnitValue(cent)
	switch currency {
	case "EUR":
		return printer.Sprintf("%.2f €", value), nil
	case "USD":
		return printer.Sprintf("$%.2f", value), nil
	default:
		meta := map[string]string{
			"currency": currency,
		}
		return "", errors.NewGenericError(500, errorDomain, errorSubDomain, ErrorUnknownCurrency, meta)
	}
}

// divide cent by 100 and return a float
func convertCentToUnitValue(cent int64) float64 {
	return float64(cent) / 100
}
