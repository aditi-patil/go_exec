package main

import (
	"net/http"
	"rest_apis/employee/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employees", api.ListAllEmployees).Methods("GET")
	router.HandleFunc("/employees", api.CreateEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", api.GetEmployee).Methods("GET")
	router.HandleFunc("/employee/{id}", api.DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/employee/{id}", api.UpdateEmployee).Methods("PUT")

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/swaggerui/").Handler(sh)

	http.ListenAndServe(":3000", router)

}
