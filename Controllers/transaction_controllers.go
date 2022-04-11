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

	idUser := "agung@mail.com"

	query := "INSERT INTO transaksi(nominal_transaksi,jenis_transaksi,tanggal_transaksi,isbn,email,kupon) VALUES (?,?,?,?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err2 := queryStatement.Exec(harga, 1, time.Now(), idBuku, idUser, -1)
	if err2 != nil {
		fmt.Println(err2)
		json.NewEncoder(w).Encode(failResponse)
		return
	}

	json.NewEncoder(w).Encode(successResponse)
}
