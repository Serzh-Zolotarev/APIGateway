package api

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := NewResponseWriterWrapper(w)
		next.ServeHTTP(ww, r)

		log.Println(time.Now().UTC().String(), "ip:", r.RemoteAddr, "Response Code:", ww.statusCode, "request_id:", r.Context().Value("request_id"))
	})
}
