package main

import (
  "fmt"
  "net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
  // text/html is the default for browsers, but Jon likes to be explicit
  // view Response Header in Web browser inspector (option+command+i)
  // under Network tab
  w.Header().Set("Content-Type", "text/html")
  // w.Header().Set("Content-Type", "text/plain")
  fmt.Println("Somebody visited our page")
  fmt.Fprint(w, "<h1>Welcome to the Awsome Sauce...</h1>")
}

func main() {
  http.HandleFunc("/", handlerFunc)
  http.ListenAndServe(":3000", nil)
}
