package DatabaseHandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type summaryInfo struct {
	client Client
	auto   Car
	order  Order
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

		query := "SELECT * FROM `order` JOIN auto ON `order`.Auto_auto_id = auto.auto_id WHERE 1=1"
		var args []interface{}

		if autoID != 0 {
			query += " AND `order`.Auto_auto_id = ?"
			args = append(args, autoID)
		}
		if clientID != 0 {
			query += " AND `order`.Client_client_id = ?"
			args = append(args, clientID)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		if rows.Next() {
			err = rows.Scan(&sum.auto.AutoId, &sum.auto.Make, &sum.auto.Model, &sum.auto.Year, &sum.auto.Vin)
			if err != nil {
				tmpl.Execute(w, "Failed to get information about car")
				return
			}
			err = rows.Scan(&sum.client.ClientId, &sum.client.Name, &sum.client.ContactInfo, &sum.client.Address)
			if err != nil {
				tmpl.Execute(w, "Failed to get information about client")
				return
			}
			//todo:more scans
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Error iterating through rows", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, sum)
	}
}
