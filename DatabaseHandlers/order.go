package DatabaseHandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Order struct {
	OrderID                 int    `json:"order_id"`
	Date                    string `json:"date"`
	Description             string `json:"description"`
	Status                  string `json:"status"`
	ClientClientID          int    `json:"Client_client_id"`
	AutoAutoID              int    `json:"Auto_auto_id"`
	ServiceStationStationID int    `json:"Service_station_station_id"`
	StaffStaffID            int    `json:"Staff_staff_id"`
	ServiceHistoryRecordID  int    `json:"Service_history_record_id"`
}

type ClientOrder struct {
	Order  Order
	Client Client
}

func GetInfoOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/order.html"))
		params := r.URL.Query().Get("Client-client-id")
		clientId, err := strconv.Atoi(params)
		if err != nil {
			tmpl.Execute(w, "")
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}

		var clientOrder ClientOrder

		query := "SELECT date,description,status FROM `order` WHERE Client_client_id = ?"
		getName := "SELECT name from client where client_id = ?"

		row := db.QueryRow(query, clientId)
		rowGetName := db.QueryRow(getName, clientId)

		err = row.Scan(&clientOrder.Order.Date, &clientOrder.Order.Description, &clientOrder.Order.Status)
		if err != nil {
			if err == sql.ErrNoRows {
				tmpl.Execute(w, "")
				http.Error(w, "Order not found", http.StatusBadRequest)
				return
			}
			tmpl.Execute(w, "")
			http.Error(w, "Failed to get information about order", http.StatusBadRequest)
			return
		}

		err = rowGetName.Scan(&clientOrder.Client.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				tmpl.Execute(w, "")
				http.Error(w, "Order not found", http.StatusBadRequest)
				return
			}
			tmpl.Execute(w, "")
			http.Error(w, "Failed to get information about order", http.StatusBadRequest)
			return
		}

		err = tmpl.Execute(w, clientOrder)
		if err != nil {
			fmt.Errorf("error executing template: %s", err)
		}
	}
}
