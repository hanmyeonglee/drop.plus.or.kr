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
		// 메인 앱: 디자인(이미지, 폰트, 미디어) 관련 외부 로딩 허용, 스크립트나 외부 폼 제출 등은 강력 차단
		csp := "default-src 'none'; script-src 'none'; style-src 'self' 'unsafe-inline'; img-src * data:; font-src * data:; media-src * data:; connect-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'none';"
		
		// 업로드된 파일(공유 링크): 인라인 렌더링 시 디자인 및 미디어 요소 허용, 스크립트 실행 등은 sandbox로 차단
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
		// Bypass authentication for static files and login/auth paths
		if strings.HasPrefix(r.URL.Path, "/static/") || strings.HasPrefix(r.URL.Path, "/auth/") || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Allow public access to GET /files/{uuid} for downloading (읽기는 누구나 가능)
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/files/") {
			next.ServeHTTP(w, r)
			return
		}

		session, _ := config.Store.Get(r, "drop-session")
		email, ok := session.Values["user_email"].(string)
		
		// If not logged in, redirect to /login
		if !ok || email == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Inject user_email into the request context
		ctx := context.WithValue(r.Context(), UserEmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
