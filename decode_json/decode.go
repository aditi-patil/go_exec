package decode_json

import (
	"encoding/json"
	"fmt"
	"log"
)

// initialise struct for books
type Book struct {
	Name      string
	Author    string
	Publisher string
	Price     float32
}

//This method decodes json and returns decoded books data
func DecodeJson(data []byte) ([]Book, error) {
	var books []Book
	err := json.Unmarshal(data, &books)
	if err != nil {
		log.Println(err)
		return books, err
	}

	fmt.Sprintf("%v", books)
	// Prints all books from decoded json
	for b := range books {
		fmt.Printf("Name: %v, Author: %v, Price: %v, Publisher: %v\n",
			books[b].Name, books[b].Author, books[b].Price, books[b].Publisher)
	}
	return books, nil
}
