package middlewares

import (
	"context"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.JWT(r).ExtractTokenMetadata(r)
		if err != nil {
			utils.Respond(w, http.StatusUnauthorized, utils.Message(err.Error()))
			return
		}

		newCtx := context.WithValue(r.Context(), utils.UserIDTag, claims.UserID)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
