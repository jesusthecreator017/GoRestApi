package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDKey contextKey = "request_id"

const RequestIDHeader string = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := req.Header.Get(RequestIDHeader)

		if id == "" {
			id = generateID()
		}

		// Store in the context for downstream handlers
		ctx := context.WithValue(req.Context(), RequestIDKey, id)
		req.WithContext(ctx)

		// Echo it back into the request header
		w.Header().Set(RequestIDHeader, id)

		next.ServeHTTP(w, req)
	})
}

func generateID() string {
	return uuid.New().String()
}

func GetRequestID(req *http.Request) string {
	id, _ := req.Context().Value(RequestIDHeader).(string)
	return id
}
