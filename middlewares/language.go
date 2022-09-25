package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
)

func LanguageParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Language")
		if lang == "" {
			lang = "EN"
		}

		c.Request = c.Request.Clone(context.WithValue(c.Request.Context(), "language", lang))

		c.Next()
	}
}
