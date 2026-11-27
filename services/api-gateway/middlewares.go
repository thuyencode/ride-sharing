package main

import (
	"net/http"
	"slices"
)

type middleware func(next http.Handler) http.HandlerFunc

func applyMiddleware(handler http.Handler, middlwares ...middleware) http.Handler {
	for _, m := range slices.Backward(middlwares) {
		handler = m(handler)
	}
	return handler
}

func prelight(origin string) middleware {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Allow prelight requests from the browser or other clients
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
