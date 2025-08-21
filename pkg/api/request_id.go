package api

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.URL.Query().Get("request_id")

		if len(requestID) == 0 {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(context.Background(), "requestID", requestID)
		log.Println("Request received with id ", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
