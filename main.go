package main

import (
	"log"
	"net/http"

	"drop.plus.or.kr/config"
	"drop.plus.or.kr/handlers"
	"drop.plus.or.kr/models"
)

func main() {
	config.LoadConfig()
	models.InitDB(config.AppConfig.DataDir)
	
	models.StartAutoDeleteScheduler()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", handlers.HandleLogin)
	mux.HandleFunc("GET /auth/login", handlers.HandleAuthLogin)
	mux.HandleFunc("GET /auth/callback", handlers.HandleAuthCallback)
	mux.HandleFunc("POST /logout", handlers.HandleLogout)
	mux.HandleFunc("GET /", handlers.HandleIndexPage)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/public/manifest.json")
	})
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/js/sw.js")
	})

	mux.HandleFunc("POST /files", handlers.HandleUploadFile)
	mux.HandleFunc("GET /files/{uuid}", handlers.HandleDownloadFile)
	mux.HandleFunc("DELETE /files/{uuid}", handlers.HandleDeleteFile)

	handler := handlers.MethodOverrideMiddleware(mux)
	handler = handlers.AuthMiddleware(handler)
	handler = handlers.SecurityHeadersMiddleware(handler)

	port := ":" + config.AppConfig.Port
	log.Printf("Server started on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
