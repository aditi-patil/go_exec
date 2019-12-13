package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"gophercises/transform/api"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	go func() {
		t := time.NewTicker(30 * time.Minute)
		for {
			select {
			case <-t.C:
				dir, _ := os.Getwd()
				imgDir := path.Join(dir, "/img")
				removeContents(imgDir)
			default:
				fmt.Print("")
			}

		}
	}()

	// registering urls with mux
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="image">
				<button type="submit">Upload Image</button>
			</form>
			</body></html>`
		fmt.Fprint(w, html)
	})
	router.HandleFunc("/modify/{id}", api.ModifyImage).Methods("GET")
	router.HandleFunc("/upload", api.UploadImage).Methods("POST")

	sh := http.StripPrefix("/img", http.FileServer(http.Dir("./img/")))
	router.PathPrefix("/img/").Handler(sh)

	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":3000", handlers.CORS(methods, origins)(router))
}

func removeContents(dir string) error {
	names, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entery := range names {
		os.RemoveAll(path.Join([]string{dir, entery.Name()}...))
	}
	return nil
}
