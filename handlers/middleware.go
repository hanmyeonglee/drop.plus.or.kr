package handlers

import (
	"context"
	"net/http"
	"strings"

	"drop.plus.or.kr/config"
)

type contextKey string

const UserEmailKey contextKey = "user_email"

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'none'; script-src 'sha256-bj3C9L6mi7diBahwpYEe91olInJomopEa6nmlb8RVto='; worker-src 'self'; style-src 'self' 'unsafe-inline'; img-src * data:; font-src * data:; media-src * data:; connect-src 'self'; manifest-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'none';"
		if strings.HasPrefix(r.URL.Path, "/files/") && r.Method == http.MethodGet {
			csp = "default-src 'none'; style-src * 'unsafe-inline'; img-src * data:; font-src * data:; media-src * data:; sandbox"
		}

		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if strings.HasPrefix(config.AppConfig.BaseURL, "https://") {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}

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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") || strings.HasPrefix(r.URL.Path, "/auth/") || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/files/") {
			next.ServeHTTP(w, r)
			return
		}

		session, _ := config.Store.Get(r, "drop-session")
		email, ok := session.Values["user_email"].(string)

		if !ok || email == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
