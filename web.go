package main 

import(
	"net/http"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
)

//Book Struct (Model)

type Book struct {
	Id 		string `json:"id"`
	Isbn 	string `json:"isbn"`
	Title 	string `json:"title"`
	Author *Author `json:"author"`
}


//Author  Struct
type Author struct{
	Firstname	string `json:"firstname"`
	Lastname	string `json:"lastname"`
}


//Init books var as a slice Book struct
var books []Book

// Get all Books
func getBooks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(books)
}


//Get single book
func getBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r) //get all params
	//Loops through books and find with id

	for _,item := range books{
		if item.Id == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// json.NewEncoder(w).Encode(&Book{})
}


//create book
func createBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = strconv.Itoa(rand.Intn(10000000))// Mock ID - not safe in production
	books = append(books,book)
	json.NewEncoder(w).Encode(book)

}


//update book
func updateBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r) //get all params
	
	for i,item := range books{
		if item.Id == params["id"]{
			books = append(books[:i], books[i+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.Id = params["id"]
			books = append(books,book)
			json.NewEncoder(w).Encode(book)
			return			
		}
	}
	json.NewEncoder(w).Encode(books)
}


//delete book
func deleteBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r) //get all params
	for i,item := range books{
		if item.Id == params["id"]{
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main(){

	//init router

	router := mux.NewRouter() //type inference

	//Mock Data @todo - implement DB
	books = append(books, Book{Id: "1", Isbn: "448742", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})

	books = append(books, Book{Id: "2", Isbn: "448744", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})
	
	//Route handlers / Endpoints

	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/book", createBook).Methods("POST")
	router.HandleFunc("/api/book/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/book/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000",router))
}