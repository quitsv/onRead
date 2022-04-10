package main

import (
	"PBPPrak/Tubes/Controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//endpoints

	router.HandleFunc("/books", Controllers.AddNewBook).Methods("POST")                                   //add new book
	router.HandleFunc("/books", Controllers.GetDetailBook).Methods("GET")                                 //get detail book
	router.HandleFunc("/books/{isbn}", Controllers.DeleteBook).Methods("DELETE")                          //delete book
	router.HandleFunc("/books/{isbn}", Controllers.UpdateBook).Methods("PUT")                             //update book
	router.HandleFunc("/pengguna", Controllers.DeleteUser).Methods("DELETE")                              //delete user
	router.HandleFunc("/books", Controllers.LookAllBookList).Methods("GET")                               // Get All List Book
	router.HandleFunc("/booksFilter/{id_genre}", Controllers.LookAllBookListFilterByGenre).Methods("GET") // Get All List Book by genre
	router.HandleFunc("/bestSeller", Controllers.LookAllBestSellerBook).Methods("GET")                    // Get All Best Seller Book
	router.HandleFunc("/bestSeller/{id_genre}", Controllers.LookAllBestSellerBookByGenre).Methods("GET")  // Get All Best Seller Book By Genre

	//cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	Handler := corsHandler.Handler(router)

	http.Handle("/", router)
	fmt.Println("Server is running on port 8080")
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", Handler))
}
