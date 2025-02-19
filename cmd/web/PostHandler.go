package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (dep *Dependencies) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
		Tmpl.ExecuteTemplate(w, "posts.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postContent := r.FormValue("postContent")
	postId := uuid.New().String()
	category := r.Form["categories"]
	title := r.FormValue("title")
	userId := r.Context().Value("user_id")

	if userId != nil {
		dep.Forum.CreatePost(postId, title, postContent, category)
	}
	fmt.Println("You have to be signed in to be able to post")
}
