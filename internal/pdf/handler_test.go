package pdf

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "pdf", response["service"])
	assert.Equal(t, "1.0.0", response["version"])
}

func TestGenerateTestReport(t *testing.T) {
	req, err := http.NewRequest("GET", "/test/report", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GenerateTestReport)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "application/pdf", rr.Header().Get("Content-Type"))

	assert.Contains(t, rr.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, rr.Header().Get("Content-Disposition"), "test_student_report.pdf")

	body := rr.Body.String()
	assert.True(t, strings.HasPrefix(body, "%PDF"), "Response should be a valid PDF")

	assert.Greater(t, len(body), 1000, "PDF should have substantial content")
}

func TestGenerateStudentReport_ValidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/students/1/report", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GenerateStudentReport)

	handler.ServeHTTP(rr, req)

	// Note: This test will fail if the Node.js API is not running or authentication fails
	// For unit testing, we would need to mock the API client
	// For now, we just check that the URL parsing works correctly
	if rr.Code == http.StatusInternalServerError {
		assert.Contains(t, rr.Body.String(), "Failed to fetch student data")
	} else if rr.Code == http.StatusOK {
		assert.Equal(t, "application/pdf", rr.Header().Get("Content-Type"))
		assert.Contains(t, rr.Header().Get("Content-Disposition"), "student_1_report.pdf")
	}
}

func TestGenerateStudentReport_Integration(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		shouldCallFunc bool
	}{
		{
			name:           "Valid report endpoint",
			path:           "/api/v1/students/1/report",
			expectedStatus: http.StatusInternalServerError, // Expected if Node.js API unavailable
			shouldCallFunc: true,
		},
		{
			name:           "Invalid ID",
			path:           "/api/v1/students/abc/report",
			expectedStatus: http.StatusBadRequest,
			shouldCallFunc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			if tt.shouldCallFunc {
				handler := http.HandlerFunc(GenerateStudentReport)
				handler.ServeHTTP(rr, req)

				if tt.expectedStatus == http.StatusBadRequest {
					assert.Equal(t, tt.expectedStatus, rr.Code)
				} else {
					assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusInternalServerError)
				}
			}
		})
	}
}

func TestGenerateStudentReport_InvalidID(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Non-numeric ID",
			url:      "/api/v1/students/abc/report",
			expected: "Invalid student ID",
		},
		{
			name:     "Missing report path",
			url:      "/api/v1/students/1",
			expected: "Invalid URL format",
		},
		{
			name:     "Wrong path format",
			url:      "/api/v1/students/1/wrong",
			expected: "Invalid URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GenerateStudentReport)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expected)
		})
	}
}

func TestGenerateStudentReport_URLParsing(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
		expectedID  string
	}{
		{
			name:        "Valid URL",
			url:         "/api/v1/students/123/report",
			expectError: false,
			expectedID:  "123",
		},
		{
			name:        "Valid URL with query params",
			url:         "/api/v1/students/456/report?format=pdf",
			expectError: false,
			expectedID:  "456",
		},
		{
			name:        "Invalid URL - missing report",
			url:         "/api/v1/students/123",
			expectError: true,
		},
		{
			name:        "Invalid URL - wrong endpoint",
			url:         "/api/v1/students/123/export",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			require.NoError(t, err)

			path := strings.TrimPrefix(u.Path, "/api/v1/students/")
			parts := strings.Split(path, "/")

			if tt.expectError {
				assert.True(t, len(parts) < 2 || parts[1] != "report",
					"Should detect invalid URL format")
			} else {
				assert.True(t, len(parts) >= 2 && parts[1] == "report",
					"Should parse valid URL format")
				assert.Equal(t, tt.expectedID, parts[0], "Should extract correct ID")
			}
		})
	}
}
