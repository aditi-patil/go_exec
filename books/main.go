package main

import (
	"books_crud_api"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	fmt.Println("Starting new api application....")

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	// registering urls with mux
	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", books_crud_api.GetBookHandler).Methods("GET")
	router.HandleFunc("/books", books_crud_api.GetAllBooksHandler).Methods("GET")
	router.HandleFunc("/book", books_crud_api.CreateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", books_crud_api.DeleteBookHandler).Methods("DELETE")
	router.HandleFunc("/book/{id}", books_crud_api.UpdateBookHandler).Methods("PUT")

	if error := http.ListenAndServe(":3000", router); error != nil {
		log.Fatal(error)
	}
}
