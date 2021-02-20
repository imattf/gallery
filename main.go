package main

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
)

// func handlerFunc(w http.ResponseWriter, r *http.Request) {
//   w.Header().Set("Content-Type", "text/html")
//   // w.Header().Set("Content-Type", "text/plain")
//
//   // logging to console on start-up
//   fmt.Println("Somebody visited our page")
//
//
//   if r.URL.Path == "/" {
//     fmt.Fprint(w, "<h1>Welcome to the Awsome Sauce...</h1>")
//   }else if r.URL.Path == "/contact" {
//     fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>.")
//   } else {
//     // have to do this first because Fprint's return StatusOk by default
//     w.WriteHeader(http.StatusNotFound)
//
//     fmt.Fprintf(w, r.URL.Path)
//
//     fmt.Fprint(w, "<h1>We could not find the page you are looking for :( </h1> <p>Please emaul us at <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a> if you keep getting sent to an invalid page.</p>")
//   }
// }

func homePage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  // print my path
  //fmt.Fprintf(w, r.URL.Path)

  fmt.Fprint(w, "<h1>Welcome to the Awesome Sauce...</h1>")

  // logging to console
  fmt.Println("home page")
}

func contactPage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  // print my path
  //fmt.Fprintf(w, r.URL.Path)

  fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>.")

  // logging to console
  fmt.Println("contact page")
}

func faqPage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  // print my path
  //fmt.Fprintf(w, r.URL.Path)

  fmt.Fprint(w, "<h1>Some Awesome FAQs...</h1>")

  // logging to console
  fmt.Println("faq page")
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", homePage)
  r.HandleFunc("/contact", contactPage)
  r.HandleFunc("/faq", faqPage)
  http.ListenAndServe(":3000", r)
}
