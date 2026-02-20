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
	json.NewDecoder(r.Body).Decode(&req)

	var id int
	var hash string
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE email=$1", req.Email).Scan(&id, &hash)
	if err != nil || utils.CheckPassword(hash, req.Password) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	token, _ := utils.GenerateJWT(id, h.Config.JWTSecret)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
