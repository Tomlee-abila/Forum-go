package main

import (
	"bytes"
	"context"
	"database/sql"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

var testDB *sql.DB

// setupTestEnvironment creates necessary test directories and files
func setupTestEnvironment() error {
	dirs := []string{
		"./ui/templates",
		"./ui/html",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	// Create minimal template files
	templates := map[string]string{
		"./ui/templates/base.html":        `{{define "base"}}{{template "content" .}}{{end}}`,
		"./ui/templates/home.html":        `{{define "content"}}{{range .}}{{.Post.Title}}{{end}}{{end}}`,
		"./ui/templates/postContent.html": `{{define "post"}}{{.Post.Title}}{{end}}`,
		"./ui/templates/categories.html":  `{{define "categories"}}{{end}}`,
		"./ui/html/error.html":            `<html><body>Error: {{.Code}}</body></html>`,
	}

	for path, content := range templates {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}

	return nil
}

// cleanupTestEnvironment removes test directories
func cleanupTestEnvironment() error {
	return os.RemoveAll("./ui")
}

func TestMain(m *testing.M) {
	// Set up test environment
	if err := setupTestEnvironment(); err != nil {
		log.Fatal(err)
	}

	// Initialize test database
	var err error
	testDB, err = models.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}
	models.DB = testDB

	// Initialize test data
	if err := initializeTestData(testDB); err != nil {
		log.Fatal(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := testDB.Close(); err != nil {
		log.Printf("Error closing test database: %v", err)
	}
	if err := cleanupTestEnvironment(); err != nil {
		log.Printf("Error cleaning up test environment: %v", err)
	}

	os.Exit(code)
}

// Helper function to set up test dependencies
func setupTestDependencies() *Dependencies {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Parse templates
	templates, err := template.ParseFiles(
		"./ui/templates/base.html",
		"./ui/templates/home.html",
		"./ui/templates/postContent.html",
		"./ui/templates/categories.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Dependencies{
		ErrorLog:  errorLog,
		InfoLog:   infoLog,
		Forum:     &models.ForumModel{DB: testDB},
		Templates: templates,
	}
}

// initializeTestData sets up test data in the database
func initializeTestData(db *sql.DB) error {
	// Drop existing tables if they exist
	_, err := db.Exec(`
		DROP TABLE IF EXISTS posts;
		DROP TABLE IF EXISTS users;
	`)
	if err != nil {
		return err
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_uuid TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			profile_picture BLOB,
			content_type TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id TEXT UNIQUE NOT NULL,
			user_uuid TEXT NOT NULL,
			username TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			media BLOB,
			content_type TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_uuid) REFERENCES users(user_uuid),
			FOREIGN KEY (username) REFERENCES users(username)
		)
	`)
	if err != nil {
		return err
	}

	// Add test user
	_, err = db.Exec(`
		INSERT INTO users (user_uuid, email, username, password)
		VALUES ('test-user-uuid', 'test@example.com', 'testuser', 'hashedpassword')
	`)
	if err != nil {
		return err
	}

	// Add test post
	_, err = db.Exec(`
		INSERT INTO posts (post_id, user_uuid, username, title, content, content_type) 
		VALUES ('test-post-id', 'test-user-uuid', 'testuser', 'Test Post', 'Test Content', '')
	`)
	return err
}

func TestHomeHandler(t *testing.T) {
	dep := setupTestDependencies()

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "Valid GET Request",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid POST Request",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			dep.HomeHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestPostHandler(t *testing.T) {
	dep := setupTestDependencies()

	tests := []struct {
		name           string
		method         string
		sessionID      string
		userUUID       string
		setupRequest   func(*http.Request)
		expectedStatus int
	}{
		{
			name:      "Valid Post Creation",
			method:    "POST",
			sessionID: "valid-session",
			userUUID:  "test-user-uuid",
			setupRequest: func(req *http.Request) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("post_content", "Test content")
				writer.WriteField("post_title", "Test title")
				writer.WriteField("categories", "Technology")
				writer.Close()

				req.Body = io.NopCloser(body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized Request",
			method:         "POST",
			sessionID:      "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/post", nil)

			if tt.setupRequest != nil {
				tt.setupRequest(req)
			}

			if tt.sessionID != "" {
				ctx := context.WithValue(req.Context(), "session_id", tt.sessionID)
				ctx = context.WithValue(ctx, "user_uuid", tt.userUUID)
				req = req.WithContext(ctx)
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: tt.sessionID,
				})
			}

			w := httptest.NewRecorder()
			dep.PostHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestLikeHandler(t *testing.T) {
	dep := setupTestDependencies()

	tests := []struct {
		name           string
		method         string
		userID         string
		formData       map[string]string
		expectedStatus int
	}{
		{
			name:   "Valid Like Request",
			method: "POST",
			userID: "test-user-id",
			formData: map[string]string{
				"id":        "test-post-id",
				"item_type": "post",
				"type":      "like",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         "GET",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.formData != nil {
				formValues := url.Values{}
				for k, v := range tt.formData {
					formValues.Set(k, v)
				}
				req = httptest.NewRequest(tt.method, "/likes", strings.NewReader(formValues.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req = httptest.NewRequest(tt.method, "/likes", nil)
			}

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "user_uuid", tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			dep.LikeHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
