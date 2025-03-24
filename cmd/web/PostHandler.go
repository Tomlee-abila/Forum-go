package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

const (
	MaxFileSize = 20 * 1024 * 1024 // 20MB to allow for some buffer
	ChunkSize   = 4096             // Read/write in 4KB chunks
)

// var DB *sql.DB

// /home/clomollo/forum/ui/html/posts.html
func (dep *Dependencies) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	userUUID, ok := r.Context().Value("user_uuid").(string)
	if !ok || userUUID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form
	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("post_title")
	content := r.FormValue("post_content")
	// categories := r.FormValue("categories")

	if title == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Create new post
	postID := uuid.New().String()
	_, err := dep.Forum.DB.Exec(`
		INSERT INTO posts (post_id, user_uuid, username, title, content)
		VALUES (?, ?, (SELECT username FROM users WHERE user_uuid = ?), ?, ?)
	`, postID, userUUID, userUUID, title, content)
	if err != nil {
		dep.ErrorLog.Printf("Failed to create post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func (dep *Dependencies) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	PostsTemplate, err := template.ParseFiles("./ui/html/posts.html")
// 	if err != nil {
// 		http.Error(w, "NOT FOUND\nError parsing post templates", http.StatusNotFound)
// 		return
// 	}
// 	posts, err := dep.Forum.AllPosts()
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Failed to get all posts", http.StatusInternalServerError)
// 		return
// 	}

// 	PostsTemplate.ExecuteTemplate(w, "allposts.html", &posts)
// }

func isValidFileType(ext string) bool {
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".mp4":  true,
		".mov":  true,
		".webm": true,
	}
	return validTypes[ext]
}

func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".webm": "video/webm",
	}
	return contentTypes[ext]
}

type CategoryFilter struct {
	Categories []string `json:"categories"`
}

func PostsByFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var filter CategoryFilter

	// Decode the JSON body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&filter); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Received Categories:", filter.Categories)

	// Call your filtering function with the extracted categories
	posts, err := models.FilterCategories(filter.Categories)
	if err != nil {
		http.Error(w, "Failed to fetch filtered posts", http.StatusInternalServerError)
		return
	}
	fmt.Println(posts)

	// Respond with the filtered posts as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
