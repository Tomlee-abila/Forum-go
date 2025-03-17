package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

type ProfilePageData struct {
	User  models.User   // current user's data
	Posts []models.Post // posts from all users
}
func (dep *Dependencies) ProfileHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("profile handler executed")

    // Check if the request is for JSON
    acceptHeader := r.Header.Get("Accept")
    isJSONRequest := acceptHeader == "application/json" || 
                     acceptHeader == "application/json, text/plain, */*" ||
                     r.Header.Get("X-Requested-With") == "XMLHttpRequest" // Ajax requests

    userID, _ := r.Context().Value("user_uuid").(string)
    user, err := dep.Forum.GetUserByID(userID)
    log.Printf("✅ user_uuid found: %s", userID)

    if err != nil || user == nil {
        log.Println("❌ user_uuid missing or not found in database!")
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // If it's a JSON request, return JSON response
    if isJSONRequest {
        log.Println("Returning JSON response for profile")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
        return
    }

    // Otherwise, serve the HTML template
    posts, err := dep.Forum.AllPosts()
    if err != nil {
        log.Printf("Error fetching posts: %v", err)
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }

    data := ProfilePageData{
        User:  *user,
        Posts: posts,
    }

    tmpl, err := template.ParseFiles("/home/clomollo/forum/ui/html/posts.html")
    if err != nil {
        http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}
