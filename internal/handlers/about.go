package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AboutHandler struct {
	DB *sql.DB
}

// PUBLIC
func (h *AboutHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var content string
	err := h.DB.QueryRow("SELECT content FROM about LIMIT 1").Scan(&content)

	if err == sql.ErrNoRows {
		// ✅ Return empty content instead of error
		json.NewEncoder(w).Encode(map[string]string{
			"content": "",
		})
		return
	}

	if err != nil {
		http.Error(w, "Failed to fetch about", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"content": content})
}

// ADMIN
func (h *AboutHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := h.DB.Exec(`
		INSERT INTO about (id, content, updated_at)
		VALUES (1, $1, NOW())
		ON CONFLICT (id)
		DO UPDATE SET content=$1, updated_at=NOW()
	`, req.Content)

	if err != nil {
		http.Error(w, "Failed to update", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}
