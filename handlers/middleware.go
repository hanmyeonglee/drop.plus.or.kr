package handlers

import (
	"net/http"
)

func MethodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			override := r.FormValue("_method")
			if override == http.MethodPut || override == http.MethodPatch || override == http.MethodDelete {
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}
