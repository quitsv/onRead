package Controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
