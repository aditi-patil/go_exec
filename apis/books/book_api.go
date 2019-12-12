// Testing go-swagger
//
// The purpose of this application is to create basic CRUD operation with mongo database, mux and go-swagger.
//
//     Schemes: http, https
//     Host: localhost:3000
//
//     Header:
//      - Access-Control-Allow-Methods: GET, POST, PUT
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// swagger:model Book
type Book struct {
	// id of this book
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	// title of this book
	Title string `bson:"title,omitempty" json:"title,omitempty"`
	// author of this book
	Author string `bson:"author,omitempty" json:"author,omitempty"`
	// publisher of this book
	Publisher string `bson:"publisher,omitempty" json:"publisher,omitempty"`
	// price of this book
	Price float32 `bson:"price,omitempty" json:"price,omitempty"`
}

var client *mongo.Client

// CreateBookHandler creates new book record into the db's book collection
func CreateBookHandler(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation POST /book CreateBookHandler
	//
	//   Create new book record into coolection db
	//   ---
	//   consumes:
	//     - application/json
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: book
	//       in: body
	//       required: true
	//       schema:
	//         type: object
	//         properties:
	//           title:
	//             type: string
	//           author:
	//             type: string
	//   responses:
	//     '200':
	//       description: book response
	//       schema:
	//         $ref: "#/definitions/Book"
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	var book Book
	_ = json.NewDecoder(request.Body).Decode(&book)
	collection := client.Database("thegoapidevelopment").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, book)
	json.NewEncoder(response).Encode(result)
}

// GetBookHandler gets book with provided objectID
func GetBookHandler(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation GET /book/{bookId} GetBookHandler
	//
	//   Returns book with specific id from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: bookId
	//       in: path
	//       description: ID of book to return
	//       required: true
	//       type: string
	//   responses:
	//     '200':
	//       description: book response
	//       schema:
	//         $ref: "#/definitions/Book"

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	collection := client.Database("thegoapidevelopment").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fmt.Println(ctx)
	err := collection.FindOne(ctx, Book{ID: id}).Decode(&book)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(book)
}

// GetAllBooksHandler gives all books from the colletion db
func GetAllBooksHandler(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation GET /books GetAllBooksHandler
	//
	//   Returns all books from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   responses:
	//     '200':
	//       description: book response
	//       schema:
	//         type: array
	//         $ref: "#/definitions/Book"
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	var books []Book
	collection := client.Database("thegoapidevelopment").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(books)
}

// DeleteBookHandler deletes data from collection with given id
func DeleteBookHandler(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation DELETE /book/{bookId} DeleteBookHandler
	//
	//   Deletes book with specific id from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: bookId
	//       in: path
	//       description: ID of book to return
	//       required: true
	//       type: string
	//   responses:
	//     '200':
	//       description: book response
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("thegoapidevelopment").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, _ := collection.DeleteOne(ctx, Book{ID: id})
	fmt.Println(reflect.TypeOf(res))
}

// UpdateBookHandler updates record
func UpdateBookHandler(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation PUT /book/{bookId} UpdateBookHandler
	//
	//   Update book parameters with specific id from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: bookId
	//       in: path
	//       description: ID of book to return
	//       required: true
	//       type: string
	//     - name: book
	//       in: body
	//       required: true
	//       schema:
	//         type: object
	//         properties:
	//           title:
	//             type: string
	//           author:
	//             type: string
	//           publisher:
	//             type: string
	//           price:
	//             type: float
	//   responses:
	//     '200':
	//       description: book updated
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	params := mux.Vars(request)
	fmt.Println(json.NewDecoder(request.Body))
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("thegoapidevelopment").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var b Book
	json.NewDecoder(request.Body).Decode(&b)
	res, _ := collection.UpdateOne(ctx, Book{ID: id}, bson.M{"$set": b})
	json.NewEncoder(response).Encode(res)
	fmt.Println(res)
}

func main() {
	fmt.Println("Starting new api application....")

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	// registering urls with mux
	router := mux.NewRouter()

	router.HandleFunc("/book/{id}", GetBookHandler).Methods("GET")
	router.HandleFunc("/books", GetAllBooksHandler).Methods("GET")
	router.HandleFunc("/book", CreateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", DeleteBookHandler).Methods("DELETE")
	router.HandleFunc("/book/{id}", UpdateBookHandler).Methods("PUT")

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/swaggerui/").Handler(sh)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router))

}
