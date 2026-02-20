package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CertificationHandler struct {
	DB *sql.DB
}

// PUBLIC
func (h *CertificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, _ := h.DB.Query("SELECT id,title,issuer,year FROM certifications ORDER BY year DESC")
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, year int
		var title, issuer string
		rows.Scan(&id, &title, &issuer, &year)
		list = append(list, map[string]interface{}{
			"id": id, "title": title, "issuer": issuer, "year": year,
		})
	}
	json.NewEncoder(w).Encode(list)
}

// ADMIN
func (h *CertificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Title  string `json:"title"`
		Issuer string `json:"issuer"`
		Year   int    `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Issuer == "" || req.Year == 0 {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec(
		"INSERT INTO certifications (title, issuer, year) VALUES ($1, $2, $3)",
		req.Title, req.Issuer, req.Year,
	)
	if err != nil {
		http.Error(w, "Create failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Certification created",
	})
}

func (h *CertificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("DELETE FROM certifications WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	// ✅ RETURN JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Certification deleted",
	})
}
