package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alsaadii98/ipctl/config"
	"github.com/alsaadii98/ipctl/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	if err := config.ConnectDB(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(config.APISecretMiddleware)

	// Routes
	r.Post("/validate", handlers.Validate)
	r.Post("/ip", handlers.AddIp)
	r.Delete("/ip", handlers.DeleteExpiredIPs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
