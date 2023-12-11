package DatabaseHandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Car struct {
	AutoId int    `json:"auto_id"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Vin    string `json:"vin"`
}

func AutoHandlerHtml(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/autopage.html"))
		params := r.URL.Query().Get("auto-id")
		autoId, err := strconv.Atoi(params)
		if err != nil {
			tmpl.Execute(w, "")
			http.Error(w, "Invalid auto ID", http.StatusBadRequest)
			return
		}

		query := "SELECT * FROM auto WHERE auto_id = ?;"
		row := db.QueryRow(query, autoId)

		var car Car

		err = row.Scan(&car.AutoId, &car.Make, &car.Model, &car.Year, &car.Vin)
		if err != nil {
			if err == sql.ErrNoRows {
				tmpl.Execute(w, "")
				http.Error(w, "Car not found", http.StatusBadRequest)
				return
			}
			tmpl.Execute(w, "")
			http.Error(w, "Failed to get information about car", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, car)
		if err != nil {
			fmt.Errorf("error executing template: %s", err)
		}
	}
}
