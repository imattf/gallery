package controllers

import (
  "gitlab.com/go-courses/lenslocked.com/models"
  "gitlab.com/go-courses/lenslocked.com/views"
  "net/http"
  "fmt"
  // "github.com/gorilla/schema"
)

type Users struct{
  NewView *views.View
  us *models.UserService
}

type SignupForm struct {
  Name     string `schema:"name"`
  Email    string `schema:"email"`
  Password string `schema:"password"`
  Age      uint   `schema:"age"`
}

func NewUsers(us *models.UserService) *Users {
  return &Users{
    NewView: views.NewView("bootstrap", "users/new"),
    us: us,
  }
}

// New is used to render the form where a user
// can create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
  if err := u.NewView.Render(w, nil); err != nil {
    panic(err)
  }
}

// Create is used to process the signup form, used
// to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
  var form SignupForm
  if err := parseForm(r, &form); err != nil {
    panic(err)
  }
  user := models.User{
    Name:  form.Name,
    Email: form.Email,
    Age:   form.Age,
  }
  if err := u.us.Create(&user); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprintln(w, form)

}
