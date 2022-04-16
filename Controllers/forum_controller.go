package Controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
)

func ViewForum(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	var forum Forum
	var arrForumResponse ArrForumResponse

	rows, err := db.Query("SELECT * FROM forum ORDER BY waktu_dikirim ASC")

	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		err := rows.Scan(&forum.IdForum, &forum.Email, &forum.WaktuDikirim, &forum.Pesan)
		if err != nil {
			fmt.Println(err)
			return
		}

		arrForumResponse.Data = append(arrForumResponse.Data, forum)
	}

	if len(arrForumResponse.Data) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arrForumResponse)
	} else {
		PrintError(400, "Tidak ada data", w)
	}
}

func WriteForum(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	_, email, _, _ := validateTokenFromCookies(r)
	waktuDikirim := time.Now()
	pesan := r.FormValue("pesan")

	_, err := db.Exec("INSERT INTO forum (email, waktu_dikirim, pesan) VALUES (?, ?, ?)", email, waktuDikirim, pesan)
	if err != nil {
		fmt.Println(err)
		PrintError(400, "Gagal menulis forum", w)
	} else {
		PrintSuccess(200, "Berhasil menulis forum", w)
	}
}
