package main

import (
  "fmt"
  "net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Somebody visited our page")
  fmt.Fprint(w, "<h1>Welcome to the Awsome Sauce...</h1>")
}

func main() {
  http.HandleFunc("/", handlerFunc)
  http.ListenAndServe(":3000", nil)
}
