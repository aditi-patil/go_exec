package main

import (
	"decode_json"
	"fmt"
)

func main() {

	jsonData := []byte(`[
			{
				"Name": "Working in Go Lang",
				"Author": "J J Roy",
				"Publisher": "TS Publications",
				"Price": 312.50
			},
			{
				"Name": "Learning in Go",
				"Author": "Thompson",
				"Publisher": "Iris Publications",
				"Price": 1230
			}
		]`)

	books, err := decode_json.DecodeJson(jsonData)
	fmt.Printf("%T\n %v, err: %v", books, books, err)

}
