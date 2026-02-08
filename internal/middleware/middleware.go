package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	APIKeyHeader = "X-API-KEY"
	ValidAPIKey  = "secret12345"
)


func APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(APIKeyHeader)
		
		if apiKey == "" || apiKey != ValidAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}


func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)


		duration := time.Since(start)
		timestamp := start.Format("2006-01-02T15:04:05")
		log.Printf("%s %s %s - Status: %d - Duration: %v",
			timestamp,
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration,
		)
	})
}


type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
