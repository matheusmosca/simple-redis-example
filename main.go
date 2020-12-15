package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	Name       string `json:"name"`
	Author     string `json:"author"`
	TotalPages int    `json:"totalpages"`
}

var books = []Book{
	{
		Name:       "And then there were none",
		Author:     "Agatha Christie",
		TotalPages: 272,
	},
	{
		Name:       "Murder on the Orient Express",
		Author:     "Agatha Christie",
		TotalPages: 256,
	},
}

func showBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")

	json.NewEncoder(w).Encode(&books)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/books", showBooks).Methods("GET")

	fmt.Println("Starting server in port 3030...")
	http.ListenAndServe(":3030", r)
}
