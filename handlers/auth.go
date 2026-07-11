package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"drop.plus.or.kr/config"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect to main if already logged in
	session, _ := config.Store.Get(r, "drop-session")
	if _, ok := session.Values["user_email"].(string); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	session, _ := config.Store.Get(r, "drop-session")
	session.Values["oauth_state"] = state
	session.Save(r, w)

	url := config.OAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "drop-session")
	savedState, ok := session.Values["oauth_state"].(string)
	
	if !ok || r.URL.Query().Get("state") != savedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := config.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("OAuth exchange failed: %v", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Extract ID token (JWT)
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token found", http.StatusInternalServerError)
		return
	}

	// Parse JWT Payload manually without verification (since it came directly from token endpoint via HTTPS)
	parts := strings.Split(rawIDToken, ".")
	if len(parts) < 2 {
		http.Error(w, "Invalid id_token format", http.StatusInternalServerError)
		return
	}
	
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		http.Error(w, "Failed to decode id_token payload", http.StatusInternalServerError)
		return
	}

	var claims struct {
		PreferredUsername string `json:"preferred_username"` // Usually the email in Entra
		Email             string `json:"email"`
	}
	json.Unmarshal(payload, &claims)

	userEmail := claims.PreferredUsername
	if userEmail == "" {
		userEmail = claims.Email
	}
	if userEmail == "" {
		http.Error(w, "Email not found in token claims", http.StatusInternalServerError)
		return
	}

	session.Values["user_email"] = userEmail
	delete(session.Values, "oauth_state")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "drop-session")
	session.Options.MaxAge = -1 // Expire cookie
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
