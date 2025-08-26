package api

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestID string
		if r.Context().Value("request_id") != nil {
			requestID = r.Context().Value("request_id").(string)
		}

		if len(requestID) == 0 {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), "request_id", requestID)
		log.Println("Request received with id ", requestID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
