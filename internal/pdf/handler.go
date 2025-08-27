package pdf

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pdf-generator/internal/api"
)

// GenerateStudentReport handles the GET /api/v1/students/:id/report endpoint
func GenerateStudentReport(w http.ResponseWriter, r *http.Request) {
	// Extract student ID from URL path: /api/v1/students/{id}/report
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/students/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "report" {
		http.Error(w, "Invalid URL format. Expected: /api/v1/students/{id}/report", http.StatusBadRequest)
		return
	}

	studentIDStr := parts[0]

	// Validate student ID
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	log.Printf("Generating PDF report for student ID: %d", studentID)

	apiClient := api.NewAPIClient()

	student, err := apiClient.GetStudentByID(studentIDStr)
	if err != nil {
		log.Printf("Error fetching student data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch student data: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched student data for: %s", student.Name)

	// Generate PDF
	pdfGenerator := NewPDFGenerator()
	pdfBytes, err := pdfGenerator.GenerateStudentReport(student)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate PDF: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully generated PDF report (%d bytes)", len(pdfBytes))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"student_%d_report.pdf\"", studentID))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	if _, err := w.Write(pdfBytes); err != nil {
		log.Printf("Error writing PDF to response: %v", err)
		return
	}

	log.Printf("PDF report successfully sent for student ID: %d", studentID)
}

// GenerateTestReport handles the test endpoint with mock data
func GenerateTestReport(w http.ResponseWriter, _ *http.Request) {
	log.Println("Generating test PDF report with mock data")

	student := api.GetMockStudent()

	pdfGenerator := NewPDFGenerator()
	pdfBytes, err := pdfGenerator.GenerateStudentReport(student)
	if err != nil {
		log.Printf("Error generating test PDF: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate test PDF: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully generated test PDF report (%d bytes)", len(pdfBytes))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"test_student_report.pdf\"")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	if _, err := w.Write(pdfBytes); err != nil {
		log.Printf("Error writing test PDF to response: %v", err)
		return
	}

	log.Println("Test PDF report successfully sent")
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "pdf",
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
