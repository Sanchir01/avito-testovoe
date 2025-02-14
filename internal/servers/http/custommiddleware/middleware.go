package custommiddleware

import (
	"context"
	"net/http"
	"strings"

	contextkey "github.com/Sanchir01/avito-testovoe/internal/context"
	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		users, err := user.ParseToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkey.UserIDCtxKey, users)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
