package middlewares

import (
	"net/http"
	"strings"

	"github.com/BenjaminRA/himnario-backend/auth"
	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/gin-gonic/gin"
)

func CheckAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/app") {
			// Check for app-specific authentication
			if err := auth.ValidateAppToken(c.Request.Header.Get("X-API-Token")); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.Next()
			return
		}

		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" || c.Request.URL.Path == "/auth/verification" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		if err := auth.VerifyToken(c); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": locale.GetLocalizedMessage(c.Request.Context().Value("language").(string), err.Error()),
			})
			return
		}

		// user, err := auth.RetrieveUser(c)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 		"error": locale.GetLocalizedMessage(c.Request.Context().Value("language").(string), err.Error()),
		// 	})
		// }
		// c.Request = c.Request.Clone(context.WithValue(c.Request.Context(), "user", user))

		c.Next()
	}
}
