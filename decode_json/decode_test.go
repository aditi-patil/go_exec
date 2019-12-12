package decode_json

import (
	"fmt"
	"testing"
)

// Test cases for DecodeJson
func TestDecodeJson(t *testing.T) {
	var jsonData []byte
	// If json is not in proper format
	jsonData = []byte(`[{"Name":"Let us C","Author":"Kanetkar", "IsPublished": "true", "Publisher": "P P Publications", "Price": 1200.00},`)

	_, error := DecodeJson(jsonData) //if error is present
	if error == nil {
		t.Error("Error is not present")
	}

	// If json data is in proper format
	jsonData = []byte(`[{"Name":"Let us C", "Author":"Kanetkar", "IsPublished": "true", "Publisher": "P P Publications", "Price": 1200.00}]`)
	books, err := DecodeJson(jsonData) //if error is not present
	if err != nil {
		t.Error("Error is present")
	}

	if books == nil { //if books are not present
		t.Errorf("Error")
	}

}

func TestDecodeJsonCheckData(t *testing.T) {
	jsonData := []byte(`[{"Name": "Test", "Author": "New Author", "Price": 100.20, "Publisher": "TT Public"}]`)
	books, error := DecodeJson(jsonData)
	fmt.Sprintf("gdfgdg %v -- %T\n%v", books, books, error)

	for b := range books {
		if books[b].Name == "" {
			t.Error("Error: Book name is not present")
		}
	}
}
