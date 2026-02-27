package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(req.Context(), duration)
			defer cancel()

			req = req.WithContext(ctx)

			// Use the custom timeoutWriter Wrapper
			tw := &timeoutWriter{
				ResponseWriter: w,
			}

			// Make the channel to handle "done" status for the ResponseWritter
			done := make(chan struct{})

			// Use concurrency to check if there has been any updates to the response writter
			go func() {
				next.ServeHTTP(tw, req)
				close(done)
			}()

			select {
			case <-done:
				// Handler finished in time
				tw.flush(w)
			case <-ctx.Done():
				if !tw.written {
					http.Error(w, "Request Timed Out", http.StatusGatewayTimeout)
				}
			}
		})
	}
}

type timeoutWriter struct {
	http.ResponseWriter
	headers http.Header
	body    []byte
	status  int
	written bool
}

func (tw *timeoutWriter) Header() http.Header {
	if tw.headers == nil {
		tw.headers = make(http.Header)
	}
	return tw.headers
}

func (tw *timeoutWriter) WriteHeader(status int) {
	tw.status = status
	tw.written = true
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.written = true
	tw.body = append(tw.body, b...)
	return len(b), nil
}

// Flush writes the buffered response to the real ResponseWriter
func (tw *timeoutWriter) flush(w http.ResponseWriter) {
	// Copy the buffered headers
	for k, v := range tw.headers {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	if tw.status != 0 {
		w.WriteHeader(tw.status)
	}

	if len(tw.body) > 0 {
		w.Write(tw.body)
	}
}
