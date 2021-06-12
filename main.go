package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/go-courses/lenslocked.com/views"
	"gitlab.com/go-courses/lenslocked.com/controllers"
	"net/http"
)

var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
)

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// have to do this first because Fprint's return StatusOk by default
	w.WriteHeader(http.StatusNotFound)

	// print my path
	//fmt.Fprintf(w, r.URL.Path)

	fmt.Fprint(w, "<h1>We could not find the page you are looking for :( </h1> <p>Please emaul us at <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a> if you keep getting sent to an invalid page.</p>")

	// logging to console
	fmt.Println("404 page")
}

func main() {

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	// instance a gorilla mux
	r := mux.NewRouter()

	// use custom 404 page
	r.NotFoundHandler = http.HandlerFunc(notFoundPage)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
