package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// LikeHandler handles likes/dislikes for both posts and comments
func (dep *Dependencies) LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	// Check user authentication
	sessionId := r.Context().Value("session_id")
	sess1, err := r.Cookie("session_id")
	if err != nil {
		log.Println("error biggy", err)
		return
	}
	if sess1.Value != sessionId {
		log.Println("sess1.Value", sess1.Value, sessionId)
		log.Println("sessioId", sessionId)
		return
	}

	userID := r.Context().Value("user_uuid").(string)

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	// Get form values
	itemID := r.FormValue("id")          // This could be either post_id or comment_id
	itemType := r.FormValue("item_type") // "post" or "comment"
	likeType := r.FormValue("type")      // "like" or "dislike"
	// Validate inputs
	if itemID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	if itemType != "post" && itemType != "comment" {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}
	if likeType != "like" && likeType != "dislike" {
		http.Error(w, "Invalid like type", http.StatusBadRequest)
		return
	}
	// Process the like/dislike based on item type
	err = dep.Forum.ProcessLike(itemType, itemID, userID, likeType)
	if err != nil {

		http.Error(w, "Failed to process like/dislike", http.StatusInternalServerError)
		return
	}

	// Get updated counts
	counts, err := dep.Forum.GetReactionCounts(itemType, itemID)
	if err != nil {
		log.Printf("Count error: %v", err)
		http.Error(w, `{"error": "Failed to get counts"}`, http.StatusInternalServerError)
		return
	}
	// Return JSON response
	response := map[string]interface{}{
		"success":  true,
		"likes":    counts["likes"],
		"dislikes": counts["dislikes"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
