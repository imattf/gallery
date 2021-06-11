package controllers

import (
  "gitlab.com/go-courses/lenslocked.com/views"
  "net/http"
  "fmt"
  // "github.com/gorilla/schema"
)

type Users struct{
  NewView *views.View
}

type SignupForm struct {
  Email    string `schema:"email"`
  Password string `schema:"password"`
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
  var form SignupForm
  if err := parseForm(r, &form); err != nil {
    panic(err)
  }

  // if err := r.ParseForm(); err != nil{
  //   panic(err)
  // }

  // dec := schema.NewDecoder()
  // var form SignupForm
  // if err:= dec.Decode(&form, r.PostForm); err != nil {
  //   panic(err)
  // }

  fmt.Fprintln(w, form)

  // r.Postform = map[string][]string
  // fmt.Fprintln(w, r.PostForm["email"])
  // fmt.Fprintln(w, r.PostForm["password"])
  // fmt.Fprintln(w, "This is a temporary response")
}
