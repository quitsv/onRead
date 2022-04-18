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

func SearchBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	var buku Buku
	var arrBuku []Buku

	isbn := r.URL.Query()["isbn"]
	judul := r.URL.Query()["judul"]
	genre := r.URL.Query()["genre"]

	queryStm := "SELECT * FROM buku"

	if len(isbn) > 0 && len(judul) > 0 && len(genre) > 0 {
		queryStm = "SELECT a.isbn, a.judul, a.penulis, a.edisi, a.tahun_cetak, a.harga, a.path_file FROM buku a JOIN genrebuku b ON a.isbn=b.isbn JOIN tipegenre c ON b.id_genre=c.id_genre WHERE b.id_genre=" + genre[0] + " AND a.isbn='" + isbn[0] + "' AND a.judul LIKE '%" + judul[0] + "%'"
	} else if len(isbn) > 0 && len(judul) > 0 {
		queryStm = "SELECT isbn, judul, penulis, edisi, tahun_cetak, harga, path_file FROM buku WHERE isbn='" + isbn[0] + "' AND judul LIKE '%" + judul[0] + "%'"
	} else if len(isbn) > 0 && len(genre) > 0 {
		queryStm = "SELECT a.isbn, a.judul, a.penulis, a.edisi, a.tahun_cetak, a.harga, a.path_file FROM buku a JOIN genrebuku b ON a.isbn=b.isbn JOIN tipegenre c ON b.id_genre=c.id_genre WHERE b.id_genre=" + genre[0] + " AND a.isbn='" + isbn[0] + "'"
	} else if len(judul) > 0 && len(genre) > 0 {
		queryStm = "SELECT a.isbn, a.judul, a.penulis, a.edisi, a.tahun_cetak, a.harga, a.path_file FROM buku a JOIN genrebuku b ON a.isbn=b.isbn JOIN tipegenre c ON b.id_genre=c.id_genre WHERE b.id_genre=" + genre[0] + " AND a.judul LIKE '%" + judul[0] + "%'"
	} else if len(isbn) > 0 {
		queryStm = "SELECT isbn, judul, penulis, edisi, tahun_cetak, harga, path_file FROM buku WHERE isbn='" + isbn[0] + "'"
	} else if len(judul) > 0 {
		queryStm = "SELECT isbn, judul, penulis, edisi, tahun_cetak, harga, path_file FROM buku WHERE judul LIKE '%" + judul[0] + "%'"
	} else if len(genre) > 0 {
		queryStm = "SELECT a.isbn, a.judul, a.penulis, a.edisi, a.tahun_cetak, a.harga, a.path_file FROM buku a JOIN genrebuku b ON a.isbn=b.isbn JOIN tipegenre c ON b.id_genre=c.id_genre WHERE b.id_genre=" + genre[0]
	}

	rows, err := db.Query(queryStm)

	if err != nil {
		fmt.Println("a:", err)
		return
	}

	for rows.Next() {
		err := rows.Scan(&buku.Isbn, &buku.Judul, &buku.Penulis, &buku.Edisi, &buku.TahunCetak, &buku.Harga, &buku.PathFile)
		if err != nil {
			fmt.Println("b:", err)
			return
		} else {
			rows2, err2 := db.Query("SELECT a.id_genre, a.genre FROM tipegenre a JOIN genrebuku b ON a.id_genre=b.id_genre WHERE b.isbn='" + buku.Isbn + "'")

			if err2 != nil {
				fmt.Println("c:", err, buku.Isbn)
				return
			}

			var genre Genre
			var arrGenre []Genre

			for rows2.Next() {
				err := rows2.Scan(&genre.IdGenre, &genre.Genre)

				if err != nil {
					fmt.Println("d:", err)
					return
				} else {
					arrGenre = append(arrGenre, genre)
				}
			}
			buku.Genre = arrGenre
			arrBuku = append(arrBuku, buku)
		}

	}

	var arrBukuResponse ArrBukuResponse
	arrBukuResponse.Data = arrBuku

	if len(arrBukuResponse.Data) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arrBukuResponse)
	} else {
		PrintError(400, "Tidak ada data", w)
	}
}
