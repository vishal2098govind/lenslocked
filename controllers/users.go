package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/vishal2098govind/lenslocked/context"
	"github.com/vishal2098govind/lenslocked/cookies"
	lenslockederr "github.com/vishal2098govind/lenslocked/errors"
	emailM "github.com/vishal2098govind/lenslocked/models/email"
	passwordResetM "github.com/vishal2098govind/lenslocked/models/password_reset"
	sessionM "github.com/vishal2098govind/lenslocked/models/session"
	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type Users struct {
	UserService *userM.UserService

	SessionService *sessionM.SessionService

	PasswordResetService *passwordResetM.PasswordResetService

	EmailService *emailM.EmailService

	Templates struct {
		New            Template // signup
		SignIn         Template // signin
		CurrentUser    Template // current user
		ForgotPassword Template // forgot password
		ResetPassword  Template // reset password
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	if context.User(r.Context()) != nil {
		http.Redirect(w, r, "/users/me", http.StatusFound)
		return
	}
	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	if context.User(r.Context()) != nil {
		http.Redirect(w, r, "/users/me", http.StatusFound)
		return
	}

	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	res, err := u.UserService.Create(userM.CreateUserRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if errors.Is(err, userM.ErrEmailAlreadyExists) {
		fmt.Println(err)
		u.Templates.New.Execute(w, r, data, lenslockederr.Public(err, "Email already exists"))
		return
	}

	if err != nil {
		fmt.Println(err)
		u.Templates.New.Execute(w, r, data, lenslockederr.Public(err, "Something went wrong"))
		return
	}

	sessRes, err := u.SessionService.Create(sessionM.CreateSessionRequest{
		UserID: res.User.ID,
	})
	if err != nil {
		fmt.Println(err)
		u.Templates.New.Execute(w, r, data, lenslockederr.Public(err, "Something went wrong"))
		return
	}

	cookies.SetCookie(w, cookies.CookieSession, sessRes.Session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	res, err := u.UserService.Authenticate(userM.AuthenticateRequest{
		Email:    email,
		Password: password,
	})
	if errors.Is(err, userM.ErrUserNotFound) {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if errors.Is(err, userM.ErrInvalidCredentials) {
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

	cookies.SetCookie(w, cookies.CookieSession, sessRes.Session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := context.Logger(ctx)
	logger.Printf("current user")

	user := context.User(ctx)

	u.Templates.CurrentUser.Execute(w, r, user)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	sessionC, err := cookies.ReadCookie(r, cookies.CookieSession)
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

	cookies.DeleteCookie(w, cookies.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	u.Templates.ForgotPassword.Execute(w, r, nil)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	user, err := u.UserService.ViaEmail(email)
	if errors.Is(err, userM.ErrUserNotFound) {
		fmt.Fprint(w, "This email is not registered. Consider creating an account using this email.")
		return
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res, err := u.PasswordResetService.Create(passwordResetM.CreatePasswordResetRequest{
		UserID: user.ID,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	val := url.Values{
		"reset-token": []string{res.PasswordReset.Token},
	}

	resetUrl := fmt.Sprintf("http://localhost:3000/reset-password?%s", val.Encode())
	err = u.EmailService.SendForgotPasswordEmail(emailM.SendForgotPasswordEmailRequest{
		To:       email,
		ResetUrl: resetUrl,
	})
	if err != nil {
		fmt.Println("failed to send email", err)
	}

	fmt.Fprintf(w, "Check your mail for reset password link")
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("reset-token")
	u.Templates.ResetPassword.Execute(w, r, struct{ Token string }{Token: token})
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("reset-token")
	newPass := r.FormValue("password")
	res, err := u.PasswordResetService.Verify(passwordResetM.VerifyRequest{
		Token: token,
	})
	if err == passwordResetM.ErrInvalidToken {
		fmt.Println(err)
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	ures, err := u.UserService.UpdatePassword(userM.UpdatePasswordRequest{
		UserID:      res.PasswordReset.UserID,
		NewPassword: newPass,
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	sess, err := u.SessionService.Create(sessionM.CreateSessionRequest{
		UserID: ures.User.ID,
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	cookies.SetCookie(w, cookies.CookieSession, sess.Session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
