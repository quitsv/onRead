package Controllers

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	email := r.Form.Get("email")
	nama := r.Form.Get("nama")
	password := r.Form.Get("password")
	tipe := 0

	result, errQuery := db.Exec("insert into pengguna (email, nama, password, tipe) values (?,?, ?, ?)", email, nama, password, tipe)

	num, _ := result.RowsAffected()

	if errQuery == nil {
		if num != 0 {
			PrintSuccess(200, "Registrasi Berhasil", w)
		} else {
			PrintError(400, "Registrasi Gagal", w)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	var pengguna Pengguna
	var arrPengguna []Pengguna

	err := r.ParseForm()
	if err != nil {
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	result, err := db.Query("select * from pengguna where tipe >= 0 and email = ? and password = ?", email, password)

	if err != nil {
		log.Print(err)
		PrintError(400, "error query", w)
	}

	for result.Next() {
		if err := result.Scan(&pengguna.Email, &pengguna.Nama, &pengguna.Password, &pengguna.Tipe); err != nil {
			log.Fatal(err.Error())
		} else {
			arrPengguna = append(arrPengguna, pengguna)
		}
	}
	if len(arrPengguna) > 0 {
		generateToken(w, arrPengguna[0].Email, arrPengguna[0].Password, arrPengguna[0].Tipe)
		PrintSuccess(200, "Login Success", w)
	} else {
		PrintError(400, "Login Failed", w)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	PrintSuccess(200, "Logged Out", w)
}
