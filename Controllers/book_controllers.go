package Controllers

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

func LookAllBookListFilterByGenre(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idGenre := vars["id_genre"]

	query := ("SELECT * from buku a join genrebuku b on a.isbn = b.isbn WHERE b.id_genre = " + idGenre)

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

func GetDetailBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	isbn := r.URL.Query()["isbn"]

	queryBuku := ("SELECT * from buku where isbn = " + isbn[0])

	rows, err := db.Query(queryBuku)
	if err != nil {
		log.Println(err)
		PrintError(400, "Error Query", w)
	}

	var book DetailBuku
	var books []DetailBuku

	for rows.Next() {
		if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
			log.Fatal(err)
			PrintError(400, "Error Fetching Data", w)
		} else {
			queryGenre := ("SELECT genrebuku.id_genre, tipegenre.genre from tipegenre join genrebuku on genrebuku.id_genre = tipegenre.id_genre where genrebuku.isbn = " + isbn[0])

			rows2, err := db.Query(queryGenre)
			if err != nil {
				log.Println(err)
				PrintError(400, "Error Query 2", w)
			}

			var genre Genre
			var genres []Genre

			for rows2.Next() {
				if err := rows2.Scan(&genre.IdGenre, &genre.Genre); err != nil {
					log.Fatal(err)
					PrintError(400, "Error Fetching Data 2", w)
				} else {
					genres = append(genres, genre)
				}
			}

			book.Genre = genres

			queryUlasan := ("select ulasan, penilaian, isbn, email from ulasanpenilaian where isbn = " + isbn[0])

			rows3, err := db.Query(queryUlasan)
			if err != nil {
				log.Println(err)
				PrintError(400, "Error Query 3", w)
			}

			var ulasan UlasanPenilaian
			var ulasans []UlasanPenilaian

			for rows3.Next() {
				if err := rows3.Scan(&ulasan.Ulasan, &ulasan.Penilaian, &ulasan.Isbn, &ulasan.Email); err != nil {
					log.Fatal(err)
					PrintError(400, "Error Fetching Data 3", w)
				} else {
					ulasans = append(ulasans, ulasan)
				}
			}
			book.Ulasan = ulasans

			books = append(books, book)
		}
	}

	var response DetailBukuResponse

	response.Data = books[0]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
