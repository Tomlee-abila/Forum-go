package main

import (
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
)

func (dep *Dependencies) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dep.HomeHandler)
	mux.Handle("/register", middleware.CSRFMiddleware(http.HandlerFunc(dep.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(dep.LogoutHandler))
	mux.Handle("/login", middleware.CSRFMiddleware(http.HandlerFunc(dep.LoginHandler)))
	return mux
}
