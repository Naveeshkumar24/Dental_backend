package routes

import (
	"github.com/gorilla/mux"

	"dental-go-backend/internal/handlers"
	"dental-go-backend/internal/middleware"
)

func Register(
	r *mux.Router,
	auth *handlers.AuthHandler,
	blog *handlers.BlogHandler,
	about *handlers.AboutHandler,
	cert *handlers.CertificationHandler,
	query *handlers.QueryHandler,
	jwtSecret string,
) {

	// ================= AUTH =================
	r.HandleFunc("/api/login", auth.Login).Methods("POST", "OPTIONS")

	// ================= PUBLIC ROUTES =================
	r.HandleFunc("/api/about", about.Get).Methods("GET")
	r.HandleFunc("/api/certifications", cert.GetAll).Methods("GET")
	r.HandleFunc("/api/blogs", blog.GetAll).Methods("GET")
	r.HandleFunc("/api/blogs/{slug}", blog.GetBySlug).Methods("GET") // ✅ FIXED
	r.HandleFunc("/api/queries", query.Create).Methods("POST")

	// ================= ADMIN ROUTES =================
	admin := r.PathPrefix("/api").Subrouter()
	admin.Use(middleware.Auth(jwtSecret))

	admin.HandleFunc("/about", about.Update).Methods("PUT")

	admin.HandleFunc("/certifications", cert.Create).Methods("POST")
	admin.HandleFunc("/certifications/{id}", cert.Delete).Methods("DELETE")

	admin.HandleFunc("/blogs", blog.Create).Methods("POST")
	admin.HandleFunc("/blogs/{id}", blog.Delete).Methods("DELETE") // ✅ slug REMOVED

	admin.HandleFunc("/queries", query.GetAll).Methods("GET")

}
