package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	UserIDKey ctxKey = "user_id"
	RoleKey   ctxKey = "role"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(UserIDKey)
	id, ok := v.(string)
	return id, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(RoleKey)
	role, ok := v.(string)
	return role, ok
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	secret := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenString := ""

			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}

			if tokenString == "" {
				c, err := r.Cookie("jwt_token")
				if err == nil {
					tokenString = c.Value
				}
			}

			if tokenString == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			type claims struct {
				UserID string `json:"user_id"`
				Role   string `json:"role"`
				jwt.RegisteredClaims
			}
			cl := &claims{}
			tkn, err := jwt.ParseWithClaims(tokenString, cl, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			})

			if err != nil || !tkn.Valid || cl.UserID == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := WithUserID(r.Context(), cl.UserID)
			ctx = WithRole(ctx, cl.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminOnly проверяет, что роль admin
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := RoleFromContext(r.Context())
		if !ok || role != "admin" {
			http.Error(w, `{"error":"admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
