package middlewares

import (
	"net/http"

	"github.com/graphql-go/handler"
)

func FinalMiddleware(next *handler.Handler) http.HandlerFunc {
	return LanguageMiddleware(next)
}
