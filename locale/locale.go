package locale

import (
	"encoding/json"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	// Initialize bundle once at startup
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("locale/en.json")
	bundle.MustLoadMessageFile("locale/es.json")
}

func GetLocalizedMessage(lang string, message_id string) string {
	// Normalize language code to lowercase
	lang = strings.ToLower(lang)
	
	// Create a new localizer for each request with the specific language
	localizer := i18n.NewLocalizer(bundle, lang)

	res, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: message_id,
	})

	if err != nil {
		return message_id
	}

	return res
}
