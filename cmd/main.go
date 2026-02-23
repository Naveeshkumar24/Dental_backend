package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"dental-go-backend/internal/config"
	"dental-go-backend/internal/database"
	"dental-go-backend/internal/handlers"
	"dental-go-backend/internal/middleware"
	"dental-go-backend/internal/routes"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// ===============================
	// CORS MIDDLEWARE
	// ===============================
	r.Use(middleware.CORS)

	// ===============================
	// 🔥 GLOBAL OPTIONS ROUTE (FIX)
	// ===============================
	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	// ===============================
	// HANDLERS
	// ===============================
	authHandler := &handlers.AuthHandler{DB: db, Config: cfg}
	blogHandler := &handlers.BlogHandler{DB: db}
	aboutHandler := &handlers.AboutHandler{DB: db}
	certHandler := &handlers.CertificationHandler{DB: db}
	queryHandler := &handlers.QueryHandler{DB: db}

	// ===============================
	// ROUTES
	// ===============================
	routes.Register(
		r,
		authHandler,
		blogHandler,
		aboutHandler,
		certHandler,
		queryHandler,
		cfg.JWTSecret,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}