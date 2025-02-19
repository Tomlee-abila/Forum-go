package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var DB *sql.DB

// function to allow the logged in user to like or dislike the posts
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	// for the users that are logged in only are the only ones to post
	cookies, err := r.Cookie("Token")
	if err != nil {
		http.Redirect(w, r, "user not logged in", http.StatusUnauthorized)
		return
	}

	userID := cookies.Value

	// Parse form data from the user/client side
	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	postID := r.FormValue("post_id")
	likeType := r.FormValue("type") // "like" or "dislike"

	// for testing use these to check the user/client vs the backend communication if it prints
	fmt.Println(postID)
	fmt.Println(likeType)

	if postID == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if the user already liked/disliked the post
	var existingType string
	err = DB.QueryRow("SELECT type FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingType)
	if err == nil {
		// Toggle the like/dislike
		if existingType != likeType {
			// if needed to for use to delet the reaction and if not use the update one(con: clicking twice)
			// 	_, err = DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			// } else {
			_, err = DB.Exec("UPDATE likes SET type = ? WHERE user_id = ? AND post_id = ?", likeType, userID, postID)
		}
	} else {
		// Insert new like/dislike
		_, err = DB.Exec("INSERT INTO likes (id, user_id, post_id, type) VALUES (?, ?, ?, ?)", uuid.New().String(), userID, postID, likeType)
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to process like/dislike", http.StatusInternalServerError)
		return
	}
	// http.Redirect(w, r, "/posts", http.StatusSeeOther) //if you want to use go for the redirecting of the page(not best for this)

	//use status code for js redirections
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("liked/disliked created successfully"))
}
