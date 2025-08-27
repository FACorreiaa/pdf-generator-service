package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DefaultNodeAPIURL = "http://localhost:5007/api/v1"
	DefaultTimeout    = 10 * time.Second
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token,omitempty"`
}

// APIClient handles communication with the Node.js backend
type APIClient struct {
	httpClient    *http.Client
	baseURL       string
	authenticated bool
	authMutex     sync.RWMutex
	authToken     string
	refreshToken  string
	csrfToken     string
}

// NewAPIClient creates a new API client
func NewAPIClient() *APIClient {
	baseURL := os.Getenv("NODE_API_URL")
	if baseURL == "" {
		baseURL = DefaultNodeAPIURL
	}

	timeout := DefaultTimeout
	if timeoutStr := os.Getenv("API_REQUEST_TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Printf("Failed to create cookie jar: %v", err)
		jar = nil
	}

	return &APIClient{
		httpClient: &http.Client{
			Timeout: timeout,
			Jar:     jar,
		},
		baseURL: baseURL,
	}
}

func (c *APIClient) authenticate() error {
	email := os.Getenv("AUTH_EMAIL")
	password := os.Getenv("AUTH_PASSWORD")

	if email == "" || password == "" {
		return fmt.Errorf("AUTH_EMAIL and AUTH_PASSWORD must be set in environment")
	}

	loginReq := LoginRequest{
		Username: email,
		Password: password,
	}

	jsonBody, err := json.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	url := fmt.Sprintf("%s/auth/login", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make login request: %w", err)
	}

	cookies := resp.Cookies()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	if loginResp.ID == 0 {
		return fmt.Errorf("login failed: invalid response")
	}

	c.authMutex.Lock()
	c.authenticated = true

	// Extract tokens from cookies - we'll send them manually
	for _, cookie := range cookies {
		if cookie.Name == "accessToken" && cookie.Value != "" {
			c.authToken = cookie.Value
			log.Printf("Stored access token: %s", c.authToken)
		}
		if cookie.Name == "refreshToken" && cookie.Value != "" {
			c.refreshToken = cookie.Value
			log.Printf("Stored refresh token: %s", c.refreshToken)
		}
		if cookie.Name == "csrfToken" && cookie.Value != "" {
			c.csrfToken = cookie.Value
			log.Printf("Stored CSRF token: %s", c.csrfToken)
		}
	}
	c.authMutex.Unlock()

	log.Printf("Successfully authenticated with Node.js backend as %s", loginResp.Name)
	return nil
}

// isAuthenticated checks if the client is authenticated
func (c *APIClient) isAuthenticated() bool {
	c.authMutex.RLock()
	defer c.authMutex.RUnlock()
	return c.authenticated
}

// ensureAuthenticated ensures the client is authenticated
func (c *APIClient) ensureAuthenticated() error {
	if !c.isAuthenticated() {
		return c.authenticate()
	}
	return nil
}

func (c *APIClient) getCsrfToken() string {
	c.authMutex.RLock()
	if c.csrfToken != "" {
		log.Printf("Using stored CSRF token: %s", c.csrfToken)
		c.authMutex.RUnlock()
		return c.csrfToken
	}
	c.authMutex.RUnlock()

	if c.httpClient.Jar == nil {
		log.Printf("No cookie jar available")
		return ""
	}

	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		log.Printf("Failed to parse base URL: %v", err)
		return ""
	}

	// ensure cookies are found
	possiblePaths := []string{"/", "/auth", "/auth/login"}
	for _, path := range possiblePaths {
		hostURL := &url.URL{
			Scheme: baseURL.Scheme,
			Host:   baseURL.Host,
			Path:   path,
		}

		cookies := c.httpClient.Jar.Cookies(hostURL)
		log.Printf("Found %d cookies for %s", len(cookies), hostURL.String())

		for _, cookie := range cookies {
			log.Printf("Cookie: %s = %s (Domain: %s, Path: %s)", cookie.Name, cookie.Value, cookie.Domain, cookie.Path)
			if cookie.Name == "csrfToken" {
				c.authMutex.Lock()
				c.csrfToken = cookie.Value
				c.authMutex.Unlock()
				log.Printf("Stored CSRF token from cookies: %s", c.csrfToken)
				return c.csrfToken
			}
		}
	}

	log.Printf("CSRF token not found in cookies")
	return ""
}

// GetStudentByID fetches student data from the Node.js API
func (c *APIClient) GetStudentByID(id string) (*Student, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	createRequest := func() (*http.Request, error) {
		url := fmt.Sprintf("%s/students/%s", c.baseURL, id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		c.authMutex.RLock()

		// Add required cookies
		var cookieParts []string
		if c.authToken != "" {
			cookieParts = append(cookieParts, fmt.Sprintf("accessToken=%s", c.authToken))
			log.Printf("Adding accessToken cookie")
		}
		if c.refreshToken != "" {
			cookieParts = append(cookieParts, fmt.Sprintf("refreshToken=%s", c.refreshToken))
			log.Printf("Adding refreshToken cookie")
		}
		if c.csrfToken != "" {
			cookieParts = append(cookieParts, fmt.Sprintf("csrfToken=%s", c.csrfToken))
			log.Printf("Adding csrfToken cookie")
		}

		if len(cookieParts) > 0 {
			cookieHeader := strings.Join(cookieParts, "; ")
			req.Header.Set("Cookie", cookieHeader)
			log.Printf("Set Cookie header with %d cookies", len(cookieParts))
		}

		// Add CSRF token as header (backend expects both cookie AND header)
		if c.csrfToken != "" {
			req.Header.Set("x-csrf-token", c.csrfToken)
			log.Printf("Adding CSRF token to request header: %s", c.csrfToken)
		} else {
			log.Printf("No CSRF token available for request")
		}

		c.authMutex.RUnlock()

		return req, nil
	}

	req, err := createRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Node.js API: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		log.Printf("Received 401, re-authenticating and retrying...")

		c.authMutex.Lock()
		c.authenticated = false
		c.authMutex.Unlock()

		if err := c.authenticate(); err != nil {
			return nil, fmt.Errorf("failed to re-authenticate after 401: %w", err)
		}

		log.Println("Retrying API request after re-authentication...")
		retryReq, err := createRequest()
		if err != nil {
			return nil, err
		}

		resp, err = c.httpClient.Do(retryReq)
		if err != nil {
			return nil, fmt.Errorf("failed to retry request after re-authentication: %w", err)
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("API request failed: %s", apiResponse.Message)
	}

	studentData, err := json.Marshal(apiResponse.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal student data: %w", err)
	}

	var student Student
	if err := json.Unmarshal(studentData, &student); err != nil {
		return nil, fmt.Errorf("failed to unmarshal student data: %w", err)
	}

	return &student, nil
}
