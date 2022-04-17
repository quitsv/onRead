package Controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func LookAllBestSellerBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	query := "SELECT b.isbn,b.judul,b.penulis,b.edisi,b.tahun_cetak,b.harga from transaksi a join buku b on a.isbn = b.isbn group by a.isbn order by count(id_transaksi) DESC;"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var book Buku
	var books []Buku

	n := 0

	for rows.Next() && n < 5 {
		if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
			log.Fatal(err)
		} else {
			books = append(books, book)
		}
		n++
	}

	var response ArrBukuResponse

	response.Data = books

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LookAllBestSellerBookByGenre(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idGenre := vars["id_genre"]

	query := "SELECT b.isbn,b.judul,b.penulis,b.edisi,b.tahun_cetak,b.harga from transaksi a join buku b on a.isbn = b.isbn join genrebuku c on c.isbn = b.isbn WHERE c.id_genre =" + idGenre + " group by a.isbn order by count(id_transaksi) DESC"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var book Buku
	var books []Buku

	n := 0

	for rows.Next() && n < 5 {
		if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
			log.Fatal(err)
		} else {
			books = append(books, book)
		}
		n++
	}

	var response ArrBukuResponse

	response.Data = books

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SearchBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	judul := r.Form.Get("judul")
	penulis := r.Form.Get("penulis")
	isbn := r.Form.Get("isbn")

	query := ("Select * from buku where judul =" + judul + " , penulis = " + penulis + " , isbn =" + isbn)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var book Buku

	if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
		log.Fatal(err)
	} else {
		book = book
	}

	var response BukuResponse

	response.Data = book

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RateBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	log.Println("BERHASIL MASUK FUNCTION DAN BERHASIL CONNECT================================================================")

	vars := mux.Vars(r)
	idBuku := vars["book_id"]

	err := r.ParseForm()
	if err != nil {
		log.Println("GAGAL DISINI")
		// log.Println(err)
		// PrintError(400, "Rating Error", w)
		return
	}

	log.Println("LEWATING PENGECEKKAN ERROR================================================================")

	ulasan := r.Form.Get("ulasan")
	penilaian, _ := strconv.Atoi(r.Form.Get("penilaian"))

	idUser := "agung@mail.com"

	log.Println("SEBELUM RESULT================================================================")

	result, errQuery := db.Exec("insert into ulasanpenilaian(ulasan,penilaian,isbn,email) values (?, ?, ?, ?)", ulasan, penilaian, idBuku, idUser)
	num, _ := result.RowsAffected()

	log.Println("SESUDAH RESULT================================================================")

	if errQuery == nil {
		if num != 0 {
			PrintSuccess(200, "Rating Given", w)
		} else {
			log.Println("GAGAL PADA PRINT ERROR================================================================")
			PrintError(400, "Failed to Rate", w)
		}
	}
}
