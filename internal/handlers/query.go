package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type QueryHandler struct {
	DB *sql.DB
}

// ==========================
// PUBLIC: CREATE QUERY
// ==========================
func (h *QueryHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.Message) == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec(
		"INSERT INTO queries (name, email, message) VALUES ($1, $2, $3)",
		req.Name,
		req.Email,
		req.Message,
	)
	if err != nil {
		http.Error(w, "Failed to submit query", http.StatusInternalServerError)
		return
	}

	// ✅ RETURN JSON RESPONSE
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Query submitted successfully",
	})
}

// ==========================
// ADMIN: GET ALL QUERIES
// ==========================
func (h *QueryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query(
		"SELECT id, name, email, message, created_at FROM queries ORDER BY created_at DESC",
	)
	if err != nil {
		http.Error(w, "Failed to fetch queries", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]map[string]interface{}, 0)

	for rows.Next() {
		var id int
		var name, email, message string
		var createdAt string

		if err := rows.Scan(&id, &name, &email, &message, &createdAt); err != nil {
			http.Error(w, "Failed to parse query", http.StatusInternalServerError)
			return
		}

		list = append(list, map[string]interface{}{
			"id":         id,
			"name":       name,
			"email":      email,
			"message":    message,
			"created_at": createdAt,
		})
	}

	json.NewEncoder(w).Encode(list)
}
