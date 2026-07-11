package main

import (
	"log"
	"net/http"

	"drop.plus.or.kr/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", handlers.HandleLogin)
	mux.HandleFunc("GET /auth/callback", handlers.HandleAuthCallback)
	mux.HandleFunc("POST /logout", handlers.HandleLogout)

	mux.HandleFunc("GET /", handlers.HandleIndexPage)

	mux.HandleFunc("POST /files", handlers.HandleUploadFile)
	mux.HandleFunc("GET /files/{uuid}", handlers.HandleDownloadFile)
	mux.HandleFunc("DELETE /files/{uuid}", handlers.HandleDeleteFile)
	
	handler := handlers.MethodOverrideMiddleware(mux)

	port := ":80"
	log.Printf("서버가 %s 포트에서 시작되었습니다.\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
