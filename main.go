package main

import (
  "fmt"
  "net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  // w.Header().Set("Content-Type", "text/plain")

  // logging to console on start-up
  fmt.Println("Somebody visited our page")


  if r.URL.Path == "/" {
    fmt.Fprint(w, "<h1>Welcome to the Awsome Sauce...</h1>")
  }else if r.URL.Path == "/contact" {
    fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>.")
  } else {
    // have to do this first because Fprint's return StatusOk by default
    w.WriteHeader(http.StatusNotFound)

    fmt.Fprintf(w, r.URL.Path)

    fmt.Fprint(w, "<h1>We could not find the page you are looking for :( </h1> <p>Please emaul us at <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a> if you keep getting sent to an invalid page.</p>")
  }
}

func main() {
  http.HandleFunc("/", handlerFunc)
  http.ListenAndServe(":3000", nil)
}
