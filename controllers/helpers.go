package controllers

import (
	"net/http"
	"net/url"

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
	return parseValues(r.PostForm, dst)
}

func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	// var form SignupForm
	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}
