package books_crud_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Router() *mux.Router {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", GetBookHandler).Methods("GET")
	router.HandleFunc("/books", GetAllBooksHandler).Methods("GET")
	router.HandleFunc("/book", CreateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", DeleteBookHandler).Methods("DELETE")
	router.HandleFunc("/book/{id}", UpdateBookHandler).Methods("PUT")

	return router
}

var result struct {
	InsertedId string
}

func TestCreateBookHandler(t *testing.T) {
	var book []byte
	book = []byte(`{"Name":"Let us C","Author":"Kanetkar", "Publisher": "P P Publications", "Price": 1200.00}`)
	request, err := http.NewRequest("POST", "/book", bytes.NewBuffer(book))
	if err != nil {
		t.Fatal("Error")
	}
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestDeleteBookHandler(t *testing.T) {
	b := []byte(addBook())
	_ = json.Unmarshal(b, &result)
	fmt.Println(result.InsertedId)

	request, err := http.NewRequest("DELETE", "/book/"+result.InsertedId, nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestUpdateBookHandler(t *testing.T) {
	var book1 []byte
	book1 = []byte(`{"Name":"Let us Go","Author":"Updated author", "Publisher": "P P Publications", "Price": 1200.00}`)
	b := []byte(addBook())
	_ = json.Unmarshal(b, &result)
	fmt.Println(result.InsertedId)
	request, err := http.NewRequest("PUT", "/book/"+result.InsertedId, bytes.NewBuffer(book1))

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestGetAllBooksHandler(t *testing.T) {
	// check response for book where id is nil
	request, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	fmt.Println(http.StatusBadRequest)

	//returns 404 Page not found error
	if status := response.Code; status == http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestGetBookHandlerForID(t *testing.T) {
	b := []byte(addBook())
	_ = json.Unmarshal(b, &result)

	// check response for book where id is not blank
	request, err := http.NewRequest("GET", "/book/"+result.InsertedId, nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	expected := `{"_id":"` + result.InsertedId + `","author":"Kanetkar","publisher":"P P Publications","price":1200}`
	if response.Body.String() != expected+"\n" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expected)
	}
}

func addBook() string {
	var book []byte
	book = []byte(`{"Name":"Let us C","Author":"Kanetkar", "Publisher": "P P Publications", "Price": 1200.00}`)
	request, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(book))
	response1 := httptest.NewRecorder()
	Router().ServeHTTP(response1, request)
	return response1.Body.String()
}
