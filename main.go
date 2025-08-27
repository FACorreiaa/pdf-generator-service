package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	pdfgen "pdf-generator/internal/pdf"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := setupRoutes()

	log.Printf("PDF Generator service starting on :%s", port)
	log.Println("Available endpoints:")
	log.Println("  GET /health - Health check")
	log.Println("  GET /test/report - Generate test PDF report with mock data")
	log.Println("  GET /api/v1/students/{id}/report - Generate student PDF report")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// setupRoutes configures all the application routes
func setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", pdfgen.HealthCheck)

	// Test endpoint for PDF generation with mock data
	mux.HandleFunc("/test/report", pdfgen.GenerateTestReport)

	// API routes handle the specific pattern for student reports
	mux.HandleFunc("/api/v1/students/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/report") {
			http.NotFound(w, r)
			return
		}
		pdfgen.GenerateStudentReport(w, r)
	})

	// Wrap with CORS middleware
	return corsMiddleware(mux)
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	allowOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*"
	}

	allowMethods := os.Getenv("CORS_ALLOWED_METHODS")
	if allowMethods == "" {
		allowMethods = "GET, POST, PUT, DELETE, OPTIONS"
	}

	allowHeaders := os.Getenv("CORS_ALLOWED_HEADERS")
	if allowHeaders == "" {
		allowHeaders = "Content-Type, Authorization"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigins)
		w.Header().Set("Access-Control-Allow-Methods", allowMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
