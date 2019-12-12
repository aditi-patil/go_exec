package main

import (
	"gophercises/stacktrace/codehandler"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", codehandler.SourceCodeHandler)
	mux.HandleFunc("/panic/", codehandler.PanicDemo)
	mux.HandleFunc("/", codehandler.Hello)
	log.Fatal(http.ListenAndServe(":3000", codehandler.DevMiddleware(mux)))
}
