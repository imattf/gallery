package controllers

import (
	"net/http"
	"time"

	"github.com/imattf/go-courses/gallery/context"
	"github.com/imattf/go-courses/gallery/email"
	"github.com/imattf/go-courses/gallery/models"
	"github.com/imattf/go-courses/gallery/rand"
	"github.com/imattf/go-courses/gallery/views"
	// "github.com/gorilla/schema"
)

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	us           models.UserService
	emailer      *email.Client
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView:      views.NewView("bootstrap", "users/new"),
		LoginView:    views.NewView("bootstrap", "users/login"),
		ForgotPwView: views.NewView("bootstrap", "users/forgot_pw"),
		ResetPwView:  views.NewView("bootstrap", "users/reset_pw"),
		us:           us,
		emailer:      emailer,
	}
}

// New is used to render the form where a user
// can create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
}

// Create is used to process the signup form, used
// to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	// this emailer could be sent on a "go" routine if needed
	u.emailer.Welcome(user.Name, user.Email)

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	alert := views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Welcome to the Gallery!",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

// Login is used to verify the provided email address and Password
// and then log the user in if they are correct
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// Logout delete the users session cookie (remember_token)
// and update the user resource with a new remember_token.
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// ResetPwForm is used to process the forgot password form
// and the reset password form.
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	// TODO: Proces the forgot password
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})

}

// ResetPw displays the reset password and has a method
// so that we can prefill the form data with a token provided
// via the URL query params.
//
// Get /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}
	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}
	u.signIn(w, user)
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
}

// Assign session cookie
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	return nil
}

// // Read email cookie for testing
// func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
// 	cookie, err := r.Cookie("remember_token")
// 	if err != nil {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}
// 	user, err := u.us.ByRemember(cookie.Value)
// 	if err != nil {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}
// 	fmt.Fprintln(w, user)
// }
