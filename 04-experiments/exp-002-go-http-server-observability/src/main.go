package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

const RequestIDKey = "request_id"

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		r = r.WithContext(context.WithValue(r.Context(), RequestIDKey, requestID))
		next.ServeHTTP(w, r)
	})
}

type Config struct {
	Port        string
	BackendAddr string
}

func loadConfig() *Config {
	return &Config{
		Port:        ":8080",
		BackendAddr: "http://localhost:8080",
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}



func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		
		next.ServeHTTP(rw, r)
		
		latency := time.Since(start)
		requestID, _ := r.Context().Value(RequestIDKey).(string)

		log.Printf(
			"method=%s path=%s status=%d latency=%v request_id=%s",
			r.Method, r.URL.Path, rw.statusCode, latency, requestID,
		)
	})
}

func main() {
	cfg := loadConfig()
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	})
	mux.HandleFunc("/unhealthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Not OK\n"))
	})

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           requestIDMiddleware(loggingMiddleware(mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Listening on %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	fmt.Println("bye")

}
