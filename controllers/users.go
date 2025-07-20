package controllers

import (
	"fmt"
	"net/http"
	"strings"

	sessionM "github.com/vishal2098govind/lenslocked/models/session"
	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type Users struct {
	UserService *userM.UserService

	SessionService *sessionM.SessionService

	Templates struct {
		New         Template // signup
		SignIn      Template // signin
		CurrentUser Template // current user
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

	sessRes, err := u.SessionService.Create(sessionM.CreateSessionRequest{
		UserID: res.User.ID,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, sessRes.Session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
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

	sessRes, err := u.SessionService.Create(sessionM.CreateSessionRequest{
		UserID: res.User.ID,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, sessRes.Session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	sessionC, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	res, err := u.SessionService.User(sessionM.GetUserIdRequest{
		Token: sessionC,
	})
	if err != nil || res.User == nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// fmt.Fprintf(w, "User: %+v\n", res.User)
	// fmt.Fprintf(w, "Headers: %+v\n", r.Header)
	u.Templates.CurrentUser.Execute(w, r, res.User)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	sessionC, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	_, err = u.SessionService.Delete(&sessionM.DeleteSessionRequest{
		Token: sessionC,
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
