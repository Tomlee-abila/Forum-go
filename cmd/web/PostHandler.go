package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
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
	postContent := r.FormValue("post_content")
	postId := uuid.New().String()
	//category := r.Form["category[]"]
	title := r.FormValue("post_title")
	user:=r.Context().Value("user_id")
	userId := r.Context().Value("user_id").(string)

	if user != nil {
		post := models.Post{
			PostId:      postId,
			UserId:      userId,
			Title:       title,
			PostContent: postContent,
		}

		dep.Forum.CreatePost(&post)

	}
	fmt.Println("You have to be signed in to be able to post")
}
