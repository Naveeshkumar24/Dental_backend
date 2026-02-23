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
	// Load environment variables (local only)
	_ = godotenv.Load()

	// Load config
	cfg := config.Load()

	// Connect database
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create router
	r := mux.NewRouter()

	// ===============================
	// CORS MIDDLEWARE (THIS IS ENOUGH)
	// ===============================
	r.Use(middleware.CORS)
	r.Use(mux.CORSMethodMiddleware(r))
	// 🔑 ALLOW OPTIONS FOR ALL ROUTES (CRITICAL FOR CORS)
	r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Method == http.MethodOptions
	}).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

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

	// ===============================
	// START SERVER
	// ===============================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
