package Controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RentBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idBuku := vars["book_id"]

	var isbn string
	var judul string
	var penulis string
	var edisi int
	var tahun_cetak time.Time
	var harga int

	var successResponse SuccessResponse
	var failResponse ErrorResponse

	successResponse.Message = "Rent Success"
	successResponse.Status = 200

	failResponse.Message = "Rent Failed"
	failResponse.Status = 400

	err := db.QueryRow("SELECT * FROM buku WHERE isbn = ?", idBuku).Scan(&isbn, &judul, &penulis, &edisi, &tahun_cetak, &harga)
	if err != nil {
		log.Println(err)
		return
	}

	_, email, _, _ := validateTokenFromCookies(r)

	query := "INSERT INTO transaksi(nominal_transaksi,jenis_transaksi,tanggal_transaksi,isbn,email,kupon) VALUES (?,?,?,?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err2 := queryStatement.Exec(harga, 1, time.Now(), idBuku, email, -1)
	if err2 != nil {
		fmt.Println(err2)
		json.NewEncoder(w).Encode(failResponse)
		return
	}

	json.NewEncoder(w).Encode(successResponse)
}

func BuyBook(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	idBuku := vars["book_id"]

	var isbn string
	var judul string
	var penulis string
	var edisi int
	var tahun_cetak time.Time
	var harga int

	var successResponse SuccessResponse
	var failResponse ErrorResponse

	successResponse.Message = "Buy Success"
	successResponse.Status = 200

	failResponse.Message = "Buy Failed"
	failResponse.Status = 400

	err := db.QueryRow("SELECT * FROM buku WHERE isbn = ?", idBuku).Scan(&isbn, &judul, &penulis, &edisi, &tahun_cetak, &harga)
	if err != nil {
		log.Println(err)
		return
	}

	_, email, _, _ := validateTokenFromCookies(r)

	query := "INSERT INTO transaksi(nominal_transaksi,jenis_transaksi,tanggal_transaksi,isbn,email,kupon) VALUES (?,?,?,?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err2 := queryStatement.Exec(harga, 2, time.Now(), idBuku, email, -1)
	if err2 != nil {
		fmt.Println(err2)
		json.NewEncoder(w).Encode(failResponse)
		return
	}

	json.NewEncoder(w).Encode(successResponse)
}

//create a function to see all transaction
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	query := "SELECT id_transaksi, nominal_transaksi, jenis_transaksi, tanggal_transaksi, isbn, email, kupon FROM transaksi"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var transactions []Transaksi

	for rows.Next() {
		var transaction Transaksi
		err = rows.Scan(&transaction.IdTransaksi, &transaction.NominalTransaksi, &transaction.JenisTransaksi, &transaction.TanggalTransaksi, &transaction.Isbn, &transaction.Email, &transaction.Kupon)
		if err != nil {
			log.Println(err)
			return
		}

		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		PrintError(400, "No Transaction Found", w)
	} else if len(transactions) == 1 {
		var response TransaksiResponse
		response.Data = transactions[0]
		json.NewEncoder(w).Encode(response)
	} else {
		var response ArrTransaksiResponse
		response.Data = transactions
		json.NewEncoder(w).Encode(response)
	}
}
