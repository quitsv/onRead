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
	Controllers.RunTools()
	router := mux.NewRouter()

	//endpoints
	router.HandleFunc("/users/login", Controllers.Login).Methods("POST")                                            // Login User
	router.HandleFunc("/users/logout", Controllers.Logout).Methods("POST")                                          // Logout User
	router.HandleFunc("/users/register", Controllers.Register).Methods("POST")                                      // Register User
	router.HandleFunc("/books", Controllers.SearchBook).Methods("GET")                                              // Get List Book by Param
	router.HandleFunc("/books/bestsellers", Controllers.LookAllBestSellerBook).Methods("GET")                       // Get All Best Seller Book
	router.HandleFunc("/books/bestsellers/{id_genre}", Controllers.LookAllBestSellerBookByGenre).Methods("GET")     // Get All Best Seller Book By Genre
	router.HandleFunc("/books/{isbn}", Controllers.GetDetailBook).Methods("GET")                                    // get detail book
	router.HandleFunc("/books/{book_id}/rating", Controllers.Authenticate(Controllers.RateBook, 0)).Methods("POST") // Rate Book (user only)
	router.HandleFunc("/books/{book_id}/rent", Controllers.Authenticate(Controllers.RentBook, 0)).Methods("POST")   // Rent Book (user only)
	router.HandleFunc("/books/{book_id}/buy", Controllers.Authenticate(Controllers.BuyBook, 0)).Methods("POST")     // Buy Book (user only)
	router.HandleFunc("/books/{isbn}/read", Controllers.Authenticate(Controllers.ReadBook, 0)).Methods("GET")       // Read Book (user only)
	router.HandleFunc("/forum", Controllers.ViewForum).Methods("GET")                                               // View Forum
	router.HandleFunc("/forum", Controllers.Authenticate(Controllers.WriteForum, 0)).Methods("POST")                // Write Forum (user only)
	router.HandleFunc("/books", Controllers.Authenticate(Controllers.AddNewBook, 1)).Methods("POST")                // add new book (admin only)
	router.HandleFunc("/books/{isbn}", Controllers.Authenticate(Controllers.DeleteBook, 1)).Methods("DELETE")       // delete book (admin only)
	router.HandleFunc("/books/{isbn}", Controllers.Authenticate(Controllers.UpdateBook, 1)).Methods("PUT")          // update book (admin only)
	router.HandleFunc("/transaksi", Controllers.Authenticate(Controllers.GetAllTransactions, 1)).Methods("GET")     // get all transactions (admin only)
	router.HandleFunc("/pengguna", Controllers.Authenticate(Controllers.DeleteUser, 1)).Methods("DELETE")           // delete user (admin only)

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
