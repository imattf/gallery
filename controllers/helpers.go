package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	// debugging the exceptions here...
	// if true {
	// 	return errors.New("blah...")
	// }
	if err := r.ParseForm(); err != nil {
		return err
	}

	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	// var form SignupForm
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
