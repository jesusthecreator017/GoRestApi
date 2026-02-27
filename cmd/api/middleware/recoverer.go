package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %s\n%s", err, debug.Stack())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
			log.Println("Middleware.Recoverer: Done")
		}()
		next.ServeHTTP(w, req)
	})
}
