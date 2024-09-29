package middlware

import (
	"context"
	"food_delivery/service"
	"net/http"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AcessTokenValdityMiddleware(next http.Handler, tokenService *service.TokenService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, err := tokenService.ValidateAccessToken(tokenService.GetTokenFromBearerString(
			r.Header.Get("Authorization")),
		)
		if err != nil {
			http.Error(w, "invalid .credentials", http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// Add claims to the request context
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}
