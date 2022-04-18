package Controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RentBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idBuku := vars["book_id"]

	idKupon := r.URL.Query()["kupon"]

	var book Buku
	var kupon Kupon

	messageSuccess := "Rent Success"
	statusSuccess := 200

	messageError := "Rent Failed"
	statusError := 400

	err := db.QueryRow("SELECT harga FROM buku WHERE isbn = ?", idBuku).Scan(&book.Harga)
	if err != nil {
		log.Println(err)
		return
	}

	errKupon := db.QueryRow("SELECT nominal,berlaku_sampai FROM kupon WHERE id_kupon = ?", idKupon[0]).Scan(&kupon.Nominal, &kupon.BerlakuSampai)
	if errKupon != nil {
		log.Println(errKupon)
		return
	}

	_, email, _, _ := validateTokenFromCookies(r)

	query := "INSERT INTO transaksi(nominal_transaksi,jenis_transaksi,tanggal_transaksi,isbn,email,kupon) VALUES (?,?,?,?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(idKupon) > 0 {
		if kupon.BerlakuSampai.Before(time.Now()) {
			_, err2 := queryStatement.Exec(book.Harga, 1, time.Now(), idBuku, email, 1)
			if err2 != nil {
				fmt.Println(err2)
				PrintError(statusError, messageError, w)
				return
			}
		} else {
			hargaKupon := book.Harga - kupon.Nominal
			_, err2 := queryStatement.Exec(hargaKupon, 1, time.Now(), idBuku, email, idKupon[0])
			if err2 != nil {
				fmt.Println(err2)
				PrintError(statusError, messageError, w)
				return
			}
		}
	}

	PrintSuccess(statusSuccess, messageSuccess, w)
	generateKupon(w, r)
}

func BuyBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idBuku := vars["book_id"]

	var book Buku
	var kupon Kupon

	idKupon := r.URL.Query()["kupon"]

	messageSuccess := "Buy Success"
	statusSuccess := 200

	messageError := "Buy Failed"
	statusError := 400

	err := db.QueryRow("SELECT harga FROM buku WHERE isbn = ?", idBuku).Scan(&book.Harga)
	if err != nil {
		log.Println(err)
		return
	}

	errKupon := db.QueryRow("SELECT nominal,berlaku_sampai FROM kupon WHERE id_kupon = ?", idKupon[0]).Scan(&kupon.Nominal, &kupon.BerlakuSampai)
	if errKupon != nil {
		log.Println(errKupon)
		return
	}

	_, email, _, _ := validateTokenFromCookies(r)

	query := "INSERT INTO transaksi(nominal_transaksi,jenis_transaksi,tanggal_transaksi,isbn,email,kupon) VALUES (?,?,?,?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(idKupon) > 0 {
		if kupon.BerlakuSampai.Before(time.Now()) {
			_, err2 := queryStatement.Exec(book.Harga, 2, time.Now(), idBuku, email, 1)
			if err2 != nil {
				fmt.Println(err2)
				PrintError(statusError, messageError, w)
				return
			}
		} else {
			hargaKupon := book.Harga - kupon.Nominal
			_, err2 := queryStatement.Exec(hargaKupon, 2, time.Now(), idBuku, email, idKupon[0])
			if err2 != nil {
				fmt.Println(err2)
				PrintError(statusError, messageError, w)
				return
			}
		}
	}

	PrintSuccess(statusSuccess, messageSuccess, w)
	generateKupon(w, r)
}

func generateKupon(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	randomNumber := rand.Intn(9)*1000 + 1000

	_, email, _, _ := validateTokenFromCookies(r)

	query := "INSERT INTO kupon(email,nominal,berlaku_sampai) VALUES (?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err2 := queryStatement.Exec(email, randomNumber, time.Now().Add(time.Hour*24*7*30))
	if err2 != nil {
		fmt.Println(err2)
		return
	}
}
