package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
	"github.com/jesusthecreator017/fswithgo/internal/auth"
)

const UserIDKey contextKey = "user_id"
const PermissionsKey contextKey = "permissions"

func GlobalAuth(jwtSecret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Get Auth header
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, req)
				return
			}

			prefix := "Bearer "
			if ok := strings.HasPrefix(authHeader, prefix); !ok {
				next.ServeHTTP(w, req)
				return
			}

			// Get and Validate the token
			token := strings.TrimPrefix(authHeader, prefix)
			id, permissions, err := auth.ValidateToken(token, jwtSecret)
			if err != nil {
				next.ServeHTTP(w, req)
				return
			}

			// Store the id and permissions in context
			ctx := context.WithValue(req.Context(), UserIDKey, id)
			ctx = context.WithValue(ctx, PermissionsKey, permissions)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func RequiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Get id from context
		_, ok := req.Context().Value(UserIDKey).(uuid.UUID)
		if !ok {
			helpers.ErrorJson(w, http.StatusUnauthorized, "id not in context")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func RequiredAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		perms := GetPermissions(req)
		if !auth.HasPermission(auth.Permission(perms), auth.PermAdmin) {
			helpers.ErrorJson(w, http.StatusForbidden, "admin access required")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func GetUserID(req *http.Request) uuid.UUID {
	id, _ := req.Context().Value(UserIDKey).(uuid.UUID)
	return id
}

func GetPermissions(req *http.Request) int32 {
	perms, _ := req.Context().Value(PermissionsKey).(int32)
	return perms
}
