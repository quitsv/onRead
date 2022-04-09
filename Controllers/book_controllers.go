package Controllers

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func LookAllBookList(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM buku"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var book Buku
	var books []Buku

	for rows.Next() {
		if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
			log.Fatal(err)
		} else {
			books = append(books, book)
		}
	}

	var response ArrBukuResponse

	response.Data = books

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
