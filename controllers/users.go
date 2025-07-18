package controllers

import (
	"fmt"
	"net/http"
	"strings"

	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type Users struct {
	UserService *userM.UserService

	Templates struct {
		New    Template // signup
		SignIn Template // signin
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	res, err := u.UserService.Create(userM.CreateUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User created: %+v", res.User)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	res, err := u.UserService.Authenticate(userM.AuthenticateRequest{
		Email:    email,
		Password: password,
	})
	if err == userM.ErrUserNotFound {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if err == userM.ErrInvalidCredentials {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "email",
		Value:    res.User.Email,
		Path:     "/",
		HttpOnly: true,
	})

	fmt.Fprint(w, "Logged user")

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	emailC, err := r.Cookie("email")
	if err != nil {
		// fmt.Fprint(w, "The email cookie could not be read")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Email cookie: %s\n", emailC.Value)
	fmt.Fprintf(w, "Headers: %+v\n", r.Header)
}
