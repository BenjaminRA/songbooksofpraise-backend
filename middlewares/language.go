package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func LanguageParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en"
		} else {
			// Normalize to lowercase and take only the language part (not region)
			lang = strings.ToLower(strings.Split(lang, "-")[0])
			
			// Validate supported languages, default to English if unsupported
			if lang != "en" && lang != "es" {
				lang = "en"
			}
		}

		c.Request = c.Request.Clone(context.WithValue(c.Request.Context(), "language", lang))

		c.Next()
	}
}
