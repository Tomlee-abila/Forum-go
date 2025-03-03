package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/utils"
)

func (dep *Dependencies) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	registerTemp, err := template.ParseFiles("./ui/html/register.html")
	if err == nil {
		if r.Method == http.MethodGet {
			csrfToken := r.Context().Value("csrf_token").(string)
			registerTemp.ExecuteTemplate(w, "register.html", map[string]interface{}{
				"CSRFToken": csrfToken,
			})
			return
		}
		if r.Method != http.MethodPost {
			dep.ClientError(w, http.StatusMethodNotAllowed)
			return
		}
		if !dep.ValidateCSRFToken(r) {
			dep.ClientError(w, http.StatusForbidden)
			return
		}

		if err := r.ParseForm(); err != nil {
			dep.ClientError(w, http.StatusBadRequest)
			return
		}

		// get the form data
		userUuid := uuid.New().String()
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
        fmt.Println(userUuid,email,username,password)
		if email == "" || username == "" || password == "" {
			dep.ClientError(w, http.StatusBadRequest)
			return
		}

		if !utils.ValidateEmail(email) {
			dep.ErrorLog.Println("Error could not validate email format")
			dep.ClientError(w, http.StatusBadRequest)
			return
		}

		userByEmail, err := dep.Forum.GetUserByEmail(email)
		if err != nil {
			dep.ServerError(w, err)
			return
		}
		if userByEmail != nil {
			dep.ClientError(w, http.StatusBadRequest)
			return
		}

		userByUsername, err := dep.Forum.GetUserByUsername(username)
		if err != nil {
			dep.ServerError(w, err)
			return
		}
		if userByUsername != nil {
			dep.ClientError(w, http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			dep.ClientError(w, http.StatusBadRequest)
			return
		}
		if err := dep.Forum.CreateUser(userUuid,email, username, password); err != nil {
			dep.ServerError(w, err)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
