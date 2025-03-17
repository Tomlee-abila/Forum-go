package main

import (
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
	// Get user session data
	userID, _ := r.Context().Value("user_uuid").(string)
	user, _ := dep.Forum.GetUserByID(userID)

	if user ==nil{
		log.Println("Error not found in database")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch posts from all users
	posts, err := dep.Forum.AllPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	// Create the data structure to pass to the template
	data := ProfilePageData{
		User:  *user,
		Posts: posts,
	}

	// Render profile template
	tmpl, err := template.ParseFiles("/home/clomollo/forum/ui/html/posts.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template, passing in the user data.
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
