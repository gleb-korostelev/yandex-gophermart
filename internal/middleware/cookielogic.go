package middleware

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/service/utils"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func EnsureUserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, err := utils.GetLoginFromCookie(r)
		if err != nil {
			logger.Infof("Failed to authorize due to error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), config.UserContextKey, login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
