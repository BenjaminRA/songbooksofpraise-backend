package locale

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle i18n.Bundle
var localizer i18n.Localizer

func GetLocalizedMessage(lang string, message_id string) string {
	// bundle := i18n.NewBundle(language.English)
	if len(bundle.LanguageTags()) == 0 {
		bundle = *i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		bundle.MustLoadMessageFile("locale/en.json")
		bundle.MustLoadMessageFile("locale/es.json")

		localizer = *i18n.NewLocalizer(&bundle, lang)
	}

	res, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: message_id,
	})

	if err != nil {
		return message_id
	}

	return res

}
