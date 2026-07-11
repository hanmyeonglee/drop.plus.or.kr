package handlers

import (
	"fmt"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Entra ID 로그인 페이지로 리다이렉트 됩니다... (구현 예정)")
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "인증 콜백 처리 및 세션 생성 (구현 예정)")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
