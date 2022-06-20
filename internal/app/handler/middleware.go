package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/romik1505/ApiGateway/internal/app/service"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func StartSpanMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.StartSpan(r.URL.Path)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func ErrorWrapper(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Printf("Error: %v", err)
			if errors.Is(err, ErrorBadRequest) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if errors.Is(err, service.ErrTokensNotFormPair) ||
				errors.Is(err, service.ErrTokensExpired) ||
				errors.Is(err, service.ErrInvalidToken) {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
