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

type Certification struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Issuer string `json:"issuer"`
	Year   int    `json:"year"`
}

// ==========================
// PUBLIC: GET ALL
// ==========================
func (h *CertificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query(
		"SELECT id, title, issuer, year FROM certifications ORDER BY year DESC",
	)
	if err != nil {
		json.NewEncoder(w).Encode([]Certification{})
		return
	}
	defer rows.Close()

	// ✅ initialize slice
	list := []Certification{}

	for rows.Next() {
		var c Certification
		if err := rows.Scan(&c.ID, &c.Title, &c.Issuer, &c.Year); err != nil {
			json.NewEncoder(w).Encode([]Certification{})
			return
		}
		list = append(list, c)
	}

	json.NewEncoder(w).Encode(list)
}

// ==========================
// ADMIN: CREATE
// ==========================
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

// ==========================
// ADMIN: DELETE
// ==========================
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

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Certification deleted",
	})
}