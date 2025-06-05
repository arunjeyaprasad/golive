package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ContextKey string

const (
	RouteParamsKey ContextKey = "routeParams"
	QueryParamsKey ContextKey = "queryParams"
)

func MuxVars(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get route parameters from mux router
		routeParams := make(map[string]string)
		vars := mux.Vars(r)
		for key, value := range vars {
			routeParams[key] = value
		}

		// Parse query parameters
		queryParams := make(map[string][]string)
		for key, values := range r.URL.Query() {
			queryParams[key] = values
		}

		// Create new context with both params
		ctx := context.WithValue(r.Context(), RouteParamsKey, routeParams)
		ctx = context.WithValue(ctx, QueryParamsKey, queryParams)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
