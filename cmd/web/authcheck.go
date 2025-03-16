package main

import (
	"encoding/json"
	"net/http"
)

func (dep *Dependencies) CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"isLoggedIn": r.Context().Value("user_uuid") != nil,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
