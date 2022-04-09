package Controllers

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func AddNewBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		PrintError(400, "Error ParseForm", w)
		return
	}

	isbn := r.Form.Get("isbn")
	judul := r.Form.Get("judul")
	penulis := r.Form.Get("penulis")
	edisi := r.Form.Get("edisi")
	tahun_cetak := r.Form.Get("tahun_cetak")
	harga := r.Form.Get("harga")

	result, errQuery := db.Exec("insert into buku (isbn, judul, penulis, edisi, tahun_cetak, harga) values (?, ?, ?, ?, ?, ?)", isbn, judul, penulis, edisi, tahun_cetak, harga)

	num, _ := result.RowsAffected()

	if errQuery == nil {
		if num != 0 {
			PrintSuccess(200, "Tambah buku berhasil", w)
		} else {
			PrintError(400, "Tambah buku gagal", w)
		}
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		PrintError(400, "Error ParseForm", w)
		return
	}

	vars := mux.Vars(r)
	isbn := vars["isbn"]

	result, errQuery := db.Exec("delete from buku where isbn = ?", isbn)

	num, _ := result.RowsAffected()

	if errQuery == nil {
		if num != 0 {
			PrintSuccess(200, "Hapus buku berhasil", w)
		} else {
			PrintError(400, "Tidak ada buku ditemukan", w)
		}
	} else {
		PrintError(400, "Hapus buku gagal", w)
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		PrintError(400, "Error ParseForm", w)
		return
	}

	vars := mux.Vars(r)
	isbn := vars["isbn"]

	isbn_new := r.Form.Get("isbn")
	judul := r.Form.Get("judul")
	penulis := r.Form.Get("penulis")
	edisi := r.Form.Get("edisi")
	tahun_cetak := r.Form.Get("tahun_cetak")
	harga := r.Form.Get("harga")

	result, errQuery := db.Exec("update buku set isbn = ?, judul = ?, penulis = ?, edisi = ?, tahun_cetak = ?, harga = ? where isbn = ?", isbn_new, judul, penulis, edisi, tahun_cetak, harga, isbn)

	num, _ := result.RowsAffected()

	if errQuery == nil {
		if num != 0 {
			PrintSuccess(200, "Update buku berhasil", w)
		} else {
			PrintError(400, "Tidak ada buku ditemukan", w)
		}
	} else {
		PrintError(400, "Update buku gagal", w)
	}
}
