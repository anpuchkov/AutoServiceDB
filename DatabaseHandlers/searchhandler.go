package DatabaseHandlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
)

type Some struct {
	Id int
	//todo
}
type summaryInfo struct {
	some Some
	//todo
}

func SearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/search.html"))

		autoIDStr := r.URL.Query().Get("auto-id")
		clientIDStr := r.URL.Query().Get("client-id")
		var sum summaryInfo
		var autoID, clientID int
		var err error

		if autoIDStr != "" {
			autoID, err = strconv.Atoi(autoIDStr)
			if err != nil {
				http.Error(w, "Ошибка преобразования параметра auto-id", http.StatusBadRequest)
				return
			}
		}

		if clientIDStr != "" {
			clientID, err = strconv.Atoi(clientIDStr)
			if err != nil {
				http.Error(w, "Ошибка преобразования параметра client-id", http.StatusBadRequest)
				return
			}
		}

		query := "SELECT * FROM auto WHERE 1=1"
		var args []interface{}

		if autoID != 0 {
			query += " AND auto_id = ?"
			args = append(args, autoID)
		}
		if clientID != 0 {
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		tmpl.Execute(w, sum)
	}
}
