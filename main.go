package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/imattf/go-courses/lenslocked.com/views"
	"net/http"
)

var (
	homeView    *views.View
	contactView *views.View
)

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}

	//Debugging stuff...
	// print my path
	// fmt.Fprintf(w, r.URL.Path)
	// logging to console
	// fmt.Println("home page")
}

func contactPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}

	//Debugging stuff...
	// print my path
	// fmt.Fprintf(w, r.URL.Path)
	// logging to console
	// fmt.Println("contact page")
}

func faqPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// print my path
	//fmt.Fprintf(w, r.URL.Path)

	fmt.Fprint(w, "<h1>Some Awesome FAQs...</h1>")

	// logging to console
	fmt.Println("faq page")
}

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

	homeView = views.NewView("views/home.gohtml")
	contactView = views.NewView("views/contact.gohtml")

	// var err error
	// homeTemplate, err = template.ParseFiles(
	//   "views/home.gohtml",
	//   "views/layouts/footer.gohtml")
	// if err != nil {
	//   panic(err)
	// }
	// contactTemplate, err = template.ParseFiles(
	//   "views/contact.gohtml",
	//   "views/layouts/footer.gohtml")
	// if err != nil {
	//   panic(err)
	// }

	// instance a gorilla mux
	r := mux.NewRouter()

	// use custom 404 page
	r.NotFoundHandler = http.HandlerFunc(notFoundPage)

	r.HandleFunc("/", homePage)
	r.HandleFunc("/contact", contactPage)
	r.HandleFunc("/faq", faqPage)
	http.ListenAndServe(":3000", r)
}
