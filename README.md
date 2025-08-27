# PDF Generator Microservice

A Go-based microservice that generates PDF reports for students by consuming the Node.js backend API.

## Features

- Consumes student data from the Node.js backend API (`/api/v1/students/:id`)
- Generates comprehensive PDF reports with student information
- Includes personal, academic, family, and address information
- RESTful API design with proper error handling
- CORS support for cross-origin requests

## Prerequisites

- Go 1.25+ installed
- PostgreSQL database with school_mgmt schema
- Node.js backend running on localhost:5007

## Installation

1. Navigate to the go-service directory:
   ```bash
   cd go-service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   go build -o pdf
   ```

## Running the Service

Start the service:
```bash
./pdf
```

The service will start on port 8080 and provide the following endpoints:

- `GET /health` - Health check endpoint
- `GET /api/v1/students/{id}/report` - Generate PDF report for student

## API Usage

### Health Check
```bash
curl http://localhost:8080/health
```

### Generate Student PDF Report
```bash
curl -o student_report.pdf http://localhost:8080/api/v1/students/1/report
```

This will download a PDF report for the student with ID 1.

## Dynamic Student ID Support

### Current Implementation Works For All Student IDs

**URL Pattern Support:**
- `/api/v1/students/1/report` → Fetches student ID **1**
- `/api/v1/students/2/report` → Fetches student ID **2**  
- `/api/v1/students/3/report` → Fetches student ID **3**
- `/api/v1/students/999/report` → Fetches student ID **999**

**How It Works:**

1. **URL Parsing** (`handler.go:17-24`):
   ```go
   path := strings.TrimPrefix(r.URL.Path, "/api/v1/students/")
   parts := strings.Split(path, "/")
   studentIDStr := parts[0]  // Extracts the ID dynamically
   ```

2. **API Call** (`handler.go:39`):
   ```go
   student, err := apiClient.GetStudentByID(studentIDStr)
   ```
   This calls: `http://localhost:5007/api/v1/students/{id}` with the dynamic ID

3. **PDF Generation**:
   Uses the fetched student data to generate a personalized PDF report

**Test It:**
```bash
# Test with different student IDs
curl -I http://localhost:8080/api/v1/students/1/report
curl -I http://localhost:8080/api/v1/students/2/report
curl -I http://localhost:8080/api/v1/students/123/report
```

## Configuration

The service is configured to connect to the Node.js API at `http://localhost:5007/api/v1`. This can be modified in the `api_client.go` file if needed.

## Architecture

The service consists of several components:

- `main.go` - HTTP server setup and routing
- `handler.go` - HTTP request handlers
- `api_client.go` - Client for communicating with Node.js API
- `student.go` - Student data structures and utility functions
- `pdf_generator.go` - PDF generation logic using gofpdf library

## Testing

To test the complete integration:

1. Ensure PostgreSQL database is running with seed data
2. Start the Node.js backend on port 5007
3. Start this Go service on port 8080
4. Use a valid student ID from the database to generate a report

Note: The Node.js API requires authentication, so you'll need valid JWT tokens to access student data.

## Testing

Run tests with:
```bash
go test ./...
```

The test suite includes:
- Health check endpoint validation
- Test PDF generation with mock data
- Student endpoint URL parsing validation
- Error handling for invalid inputs

## Docker Deployment

### Build Docker Image
```bash
docker build -t pdf:latest .
```

### Run with Docker
```bash
docker run -p 8080:8080 pdf:latest
```

### Run with Docker Compose
```bash
docker-compose up
```

The Docker setup includes:
- Multi-stage build for smaller image size
- Health checks for monitoring
- Environment variable configuration
- Alpine Linux base for security

## Observability

The service uses Prometheus for monitoring and Grafana for visualization.
Not completed because its not requited on the project requirements, but could have easily have set grafana and all the
services to run locally from the compose and see all metrics on Grafana dashboard.

## Dependencies

- `github.com/jung-kurt/gofpdf` - PDF generation library
- `github.com/stretchr/testify` - Testing framework

## Architecture Changes

The service has been refactored to:
- ✅ Use Go's standard HTTP router
- ✅ Include comprehensive test coverage with testify
- ✅ Support Docker containerization
- ✅ Environment-based configuration
- ✅ Multi-stage Docker builds for optimized images
- ✅ Mock observability for testing
