package books_crud_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Book struct details
type Book struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Author    string             `bson:"author,omitempty" json:"author,omitempty"`
	Publisher string             `bson:"publisher,omitempty" json:"publisher,omitempty"`
	Price     float32            `bson:"price,omitempty" json:"price,omitempty"`
}

var client *mongo.Client

// CreateBookHandler creates new book record into the db's book collection
func CreateBookHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
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
	response.Header().Set("content-type", "application/json")
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
	response.Header().Set("content-type", "application/json")
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
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	fmt.Println(json.NewDecoder(request.Body))
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("thegoapidevelopment").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var b Book
	json.NewDecoder(request.Body).Decode(&b)
	res, _ := collection.UpdateOne(ctx, Book{ID: id}, bson.M{"$set": b})
	fmt.Println(res)

}
