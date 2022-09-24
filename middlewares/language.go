package middlewares

import (
	"context"
	"net/http"

	"github.com/graphql-go/handler"
)

func LanguageMiddleware(next *handler.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Language")
		if lang == "" {
			lang = "EN"
		}
		ctx := context.WithValue(r.Context(), "language", lang)

		next.ContextHandler(ctx, w, r)
	})
}
