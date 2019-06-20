package main

import (
  "encoding/json"
  "log"
  "net/http"
  "os"

  "github.com/joho/godotenv"
  "github.com/julienschmidt/httprouter"
)

var listen string

func main() {
  // Load .env
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  listen = os.Getenv("LISTEN")

  // Routes
  router := httprouter.New()

  // Start server
  log.Printf("starting server on %s", listen)
  log.Fatal(http.ListenAndServe(listen, router))
}
