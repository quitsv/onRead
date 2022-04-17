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

	json.NewEncoder(w).Encode(arrForumResponse)
}

func WriteForum(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	var arrForumResponse ArrForumResponse

	email := r.FormValue("email")
	waktuDikirim := time.Now()
	pesan := r.FormValue("pesan")

	_, err := db.Exec("INSERT INTO forum (email, waktu_dikirim, pesan) VALUES (?, ?, ?)", email, waktuDikirim, pesan)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(arrForumResponse)
}
