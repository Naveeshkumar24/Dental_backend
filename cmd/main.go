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
	// Load environment variables
	_ = godotenv.Load()
	cfg := config.Load()

	// Connect database
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Router
	r := mux.NewRouter()

	// ✅ REQUIRED FOR CORS + OPTIONS
	r.Use(middleware.CORS)
	r.Use(mux.CORSMethodMiddleware(r))

	// Handlers
	authHandler := &handlers.AuthHandler{DB: db, Config: cfg}
	blogHandler := &handlers.BlogHandler{DB: db}
	aboutHandler := &handlers.AboutHandler{DB: db}
	certHandler := &handlers.CertificationHandler{DB: db}
	queryHandler := &handlers.QueryHandler{DB: db}

	// Routes
	routes.Register(
		r,
		authHandler,
		blogHandler,
		aboutHandler,
		certHandler,
		queryHandler,
		cfg.JWTSecret,
	)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local fallback
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
