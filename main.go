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
    fmt.Fprintf(w, r.URL.Path)
  }
}

func main() {
  http.HandleFunc("/", handlerFunc)
  http.ListenAndServe(":3000", nil)
}
