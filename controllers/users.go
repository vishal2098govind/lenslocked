package controllers

import (
	"fmt"
	"net/http"

	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type Users struct {
	UserService *userM.UserService

	Templates struct {
		New Template // signup
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, data)
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
