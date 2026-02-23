package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"dental-go-backend/internal/utils"
	"dental-go-backend/internal/config"
)

type AuthHandler struct {
	DB     *sql.DB
	Config *config.Config
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string
		Password string
	}

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var id int
	var hash string

	err := h.DB.QueryRow(
		"SELECT id, password FROM users WHERE email=$1",
		req.Email,
	).Scan(&id, &hash)

	if err != nil || utils.CheckPassword(hash, req.Password) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 🔴 DO NOT IGNORE THIS ERROR
	token, err := utils.GenerateJWT(id, h.Config.JWTSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// ✅ ALWAYS SEND A PROPER RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}