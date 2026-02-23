package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type BlogHandler struct {
	DB *sql.DB
}

type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ==========================
// GET ALL BLOGS (PUBLIC)
// ==========================
func (h *BlogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query(`
		SELECT id, title, slug, created_at
		FROM blogs
		ORDER BY created_at DESC
	`)
	if err != nil {
		json.NewEncoder(w).Encode([]Blog{})
		return
	}
	defer rows.Close()

	// ✅ IMPORTANT: initialize empty slice
	blogs := []Blog{}

	for rows.Next() {
		var b Blog
		if err := rows.Scan(&b.ID, &b.Title, &b.Slug, &b.CreatedAt); err != nil {
			json.NewEncoder(w).Encode([]Blog{})
			return
		}
		blogs = append(blogs, b)
	}

	// ✅ check iteration error
	if err := rows.Err(); err != nil {
		json.NewEncoder(w).Encode([]Blog{})
		return
	}

	json.NewEncoder(w).Encode(blogs)
}

// ==========================
// GET BLOG BY SLUG (PUBLIC)
// ==========================
func (h *BlogHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := mux.Vars(r)["slug"]

	var blog Blog

	err := h.DB.QueryRow(`
		SELECT id, title, slug, content, created_at
		FROM blogs
		WHERE slug = $1
	`, slug).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Slug,
		&blog.Content,
		&blog.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch blog", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(blog)
}

// ==========================
// CREATE BLOG (ADMIN)
// ==========================
func (h *BlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Title), " ", "-"))

	var blog Blog

	err := h.DB.QueryRow(`
		INSERT INTO blogs (title, slug, content)
		VALUES ($1, $2, $3)
		RETURNING id, title, slug, content, created_at
	`, req.Title, slug, req.Content).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Slug,
		&blog.Content,
		&blog.CreatedAt,
	)

	if err != nil {
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

// ==========================
// DELETE BLOG (ADMIN)
// ==========================
func (h *BlogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog deleted",
	})
}
