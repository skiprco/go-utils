package converters

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// return a printer with the expected language. If the language is not found or empty, return default language (English)
func getPrinter(expectedLang string) *message.Printer {
	lang := language.English
	if expectedLang != "" {
		lang = language.Make(expectedLang)
	}
	return message.NewPrinter(lang)
}
