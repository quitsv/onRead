package Controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func LookAllBookList(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM buku"

	isbn := r.URL.Query()["isbn"]
	judul := r.URL.Query()["judul"]

	if len(isbn) > 0 && len(judul) > 0 {
		query = query + " WHERE isbn = " + isbn[0] + " AND judul = '" + judul[0] + "'"
	} else if len(isbn) > 0 {
		query = query + " WHERE isbn = " + isbn[0]
	} else if len(judul) > 0 {
		query = query + " WHERE judul like '%" + judul[0] + "%'"
	}

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
			queryGenre := ("SELECT genrebuku.id_genre, tipegenre.genre from tipegenre join genrebuku on genrebuku.id_genre = tipegenre.id_genre where isbn = " + book.Isbn)

			rows2, err := db.Query(queryGenre)
			if err != nil {
				log.Println(err)
				PrintError(400, "Error Query Genre", w)
			}

			var genre Genre
			var genres []Genre

			for rows2.Next() {
				if err := rows2.Scan(&genre.IdGenre, &genre.Genre); err != nil {
					log.Fatal(err)
					PrintError(400, "Error Fetching Genre", w)
				} else {
					genres = append(genres, genre)
				}
			}

			book.Genre = genres
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

	idGenre := r.URL.Query()["id_genre"]

	query := ("SELECT * from buku a join genrebuku b on b.isbn = a.isbn")

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
			queryGenre := ("SELECT genrebuku.id_genre, tipegenre.genre from tipegenre join genrebuku on genrebuku.id_genre = tipegenre.id_genre where genrebuku.id_genre = " + idGenre[0])

			rows2, err := db.Query(queryGenre)
			if err != nil {
				log.Println(err)
				PrintError(400, "Error Query Genre", w)
			}

			var genre Genre
			var genres []Genre

			for rows2.Next() {
				if err := rows2.Scan(&genre.IdGenre, &genre.Genre); err != nil {
					log.Fatal(err)
					PrintError(400, "Error Fetching Genre", w)
				} else {
					genres = append(genres, genre)
				}
			}

			book.Genre = genres
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

	vars := mux.Vars(r)

	isbn := vars["isbn"]

	queryBuku := ("SELECT * from buku where isbn = " + isbn)

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
			queryGenre := ("SELECT genrebuku.id_genre, tipegenre.genre from tipegenre join genrebuku on genrebuku.id_genre = tipegenre.id_genre where genrebuku.isbn = " + isbn)

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

			queryUlasan := ("select ulasan, penilaian, isbn, email from ulasanpenilaian where isbn = " + isbn)

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

	idGenre := r.URL.Query()["genre"]

	query := "SELECT * from view_penjualan_buku_per_genre "

	if len(idGenre) > 0 {
		query = query + " where genre = '" + idGenre[0] + "'"
	}
	query = query + "order by jumlah_penjualan DESC"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var bestSeller BestSeller
	var bestSellers []BestSeller

	for rows.Next() {
		if err := rows.Scan(&bestSeller.Isbn, &bestSeller.Judul, &bestSeller.Id_genre, &bestSeller.Genre, &bestSeller.Jumlah_penjualan); err != nil {
			log.Fatal(err)
			PrintError(400, "Error Fetching Data", w)
		} else {
			queryBuku := ("Select buku.isbn,buku.judul,buku.penulis,buku.edisi,buku.tahun_cetak,buku.harga from buku where isbn = '" + bestSeller.Isbn + "'")

			rows2, err := db.Query(queryBuku)
			if err != nil {
				log.Println(err)
				PrintError(400, "Error Query Book", w)
			}

			var book Buku
			var books []Buku

			for rows2.Next() {
				if err := rows2.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
					log.Fatal(err)
					PrintError(400, "Error Fetching Book", w)
				} else {
					books = append(books, book)

				}
			}

			bestSeller.Buku = books

			queryGenre := ("Select tipegenre.id_genre,tipegenre.genre from tipegenre join genrebuku on tipegenre.id_genre = genrebuku.id_genre where genrebuku.isbn ='" + bestSeller.Isbn + "'")
			rows3, err := db.Query(queryGenre)
			if err != nil {
				log.Println(err)
				PrintError(400, "Error Query Book", w)
			}

			var genre Genre
			var genres []Genre

			for rows3.Next() {
				if err := rows3.Scan(&genre.IdGenre, &genre.Genre); err != nil {
					log.Fatal(err)
					PrintError(400, "Error Fetching Book", w)
				} else {
					genres = append(genres, genre)

				}
			}
			bestSeller.Genres = genres
			bestSellers = append(bestSellers, bestSeller)
		}
	}

	var response BestSellerResponse

	response.Data = bestSellers

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func LookAllBestSellerBookByGenre(w http.ResponseWriter, r *http.Request) {
// 	db := Connect()
// 	defer db.Close()

// 	vars := mux.Vars(r)
// 	idGenre := vars["id_genre"]

// 	query := "SELECT b.isbn,b.judul,b.penulis,b.edisi,b.tahun_cetak,b.harga from transaksi a join buku b on a.isbn = b.isbn join genrebuku c on c.isbn = b.isbn WHERE c.id_genre =" + idGenre + " group by a.isbn order by count(id_transaksi) DESC"

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var book Buku
// 	var books []Buku

// 	n := 0

// 	for rows.Next() && n < 5 {
// 		if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
// 			log.Fatal(err)
// 		} else {
// 			books = append(books, book)
// 		}
// 		n++
// 	}

// 	var response ArrBukuResponse

// 	response.Data = books

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func SearchBook(w http.ResponseWriter, r *http.Request) {
// 	db := Connect()
// 	defer db.Close()

// 	judul := r.Form.Get("judul")
// 	penulis := r.Form.Get("penulis")
// 	isbn := r.Form.Get("isbn")

// 	query := ("Select * from buku where judul =" + judul + " , penulis = " + penulis + " , isbn =" + isbn)

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var book Buku

// 	if err := rows.Scan(&book.Isbn, &book.Judul, &book.Penulis, &book.Edisi, &book.TahunCetak, &book.Harga); err != nil {
// 		log.Fatal(err)
// 	}

// 	var response BukuResponse

// 	response.Data = book

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func ReadBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	_, email, _, _ := validateTokenFromCookies(r)

	vars := mux.Vars(r)
	isbn := vars["isbn"]

	var transaksi Transaksi
	var arrTransaksi []Transaksi

	rows, err := db.Query("SELECT id_transaksi, jenis_transaksi, tanggal_transaksi, isbn, email FROM transaksi WHERE email = ? AND isbn = ?", email, isbn)

	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		err := rows.Scan(&transaksi.IdTransaksi, &transaksi.JenisTransaksi, &transaksi.TanggalTransaksi, &transaksi.Isbn, &transaksi.Email)
		if err != nil {
			fmt.Println(err)
			return
		}

		arrTransaksi = append(arrTransaksi, transaksi)
	}

	haveBook := false
	for i := 0; i < len(arrTransaksi); i++ {
		if arrTransaksi[i].JenisTransaksi == "beli" {
			haveBook = true
			break
		} else if arrTransaksi[i].JenisTransaksi == "sewa" && arrTransaksi[i].TanggalTransaksi.AddDate(0, 0, 30).After(time.Now()) {
			haveBook = true
			break
		}
	}

	if haveBook {
		rows, err := db.Query("SELECT * FROM buku WHERE isbn = ?", isbn)

		if err != nil {
			fmt.Println(err)
			PrintError(400, "error query", w)
		}

		var buku Buku
		var arrBuku []Buku
		var bukuResponse BukuResponse

		for rows.Next() {
			if err := rows.Scan(&buku.Isbn, &buku.Judul, &buku.Penulis, &buku.Edisi, &buku.TahunCetak, &buku.Harga, &buku.PathFile); err != nil {
				fmt.Println("err 2")
				fmt.Println(err)
			} else {
				arrBuku = append(arrBuku, buku)
			}
		}

		if len(arrBuku) > 0 {
			bukuResponse.Data = arrBuku[0]
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bukuResponse)
		} else {
			PrintError(400, "Tidak ada data", w)
		}
	} else {
		PrintError(404, "Anda tidak memiliki buku tersebut", w)
	}
}

func RateBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idBuku := vars["book_id"]

	err := r.ParseForm()
	if err != nil {
		log.Println("GAGAL DISINI")
		// log.Println(err)
		// PrintError(400, "Rating Error", w)
		return
	}

	ulasan := r.Form.Get("ulasan")
	penilaian, _ := strconv.Atoi(r.Form.Get("penilaian"))

	idUser := "agung@mail.com"

	result, errQuery := db.Exec("insert into ulasanpenilaian(ulasan,penilaian,isbn,email) values (?, ?, ?, ?)", ulasan, penilaian, idBuku, idUser)
	num, _ := result.RowsAffected()

	if errQuery == nil {
		if num != 0 {
			PrintSuccess(200, "Rating Given", w)
		} else {
			PrintError(400, "Failed to Rate", w)
		}
	}
}
