package controllers

import (
  "gitlab.com/go-courses/lenslocked.com/views"
  "net/http"
  "fmt"
)

type Users struct{
  NewView *views.View
}

func NewUsers() *Users {
  return &Users{
    NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
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
  fmt.Fprintln(w, "This is a temporary response")
}
