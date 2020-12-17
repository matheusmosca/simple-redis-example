package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

func newRedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return rdb
}

type Book struct {
	Name       string `json:"name"`
	Author     string `json:"author"`
	TotalPages int    `json:"totalpages"`
}

//* Repository

func getBook(bookName string) (*Book, error) {
	rdb := newRedisConnection()
	defer rdb.Close()

	bookJSON, err := rdb.Get(ctx, bookName).Result()
	if err != nil {
		return nil, err
	}

	rawBook := &Book{}
	err = json.Unmarshal([]byte(bookJSON), rawBook)
	if err != nil {
		return nil, err
	}

	return rawBook, nil
}

func storeBook(book *Book) error {
	rdb := newRedisConnection()
	defer rdb.Close()

	_, err := getBook(book.Name)
	if err == nil {
		return errors.New("This key name already exists on the database")
	}

	bookJSON, err := json.Marshal(book)

	if err != nil {
		return err
	}

	if err = rdb.Set(ctx, book.Name, bookJSON, 0).Err(); err != nil {
		return err
	}

	return nil
}

//* Handlers

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")
	newBook := &Book{}

	err := json.NewDecoder(r.Body).Decode(newBook)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = storeBook(newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(newBook)
}

func showBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")

	params := mux.Vars(r)

	book, err := getBook(params["name"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{name}", showBook).Methods("GET")

	fmt.Println("Starting server in port 3030...")
	http.ListenAndServe(":3030", r)
}
