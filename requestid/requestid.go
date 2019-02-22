package requestid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	requestIDContextKey = contextKey("X-Request-ID")
)

// Middleware extracts the request ID from a request's headers, if it is
// present, and stores it in the request's context.
//
// If no request ID is present in the request's headers, it will be generated.
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext extracts the request ID from a request's context.
func FromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("requestid: context is nil")
	}

	requestID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return "", fmt.Errorf("requestid: %s not found in context", requestIDContextKey)
	}

	return requestID, nil
}
