package DatabaseHandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

type Client struct {
	ClientId    int    `json:"client_id"`
	Name        string `json:"name"`
	ContactInfo string `json:"contact_info"`
	Address     string `json:"address"`
}

func ClientHandlerHtml(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/clientpage.html"))
		params := r.URL.Query().Get("client-id")
		clientId, err := strconv.Atoi(params)
		if err != nil {
			tmpl.Execute(w, "")
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}

		query := "SELECT * FROM client WHERE client_id = ?"
		row := db.QueryRow(query, clientId)

		var client Client

		err = row.Scan(&client.ClientId, &client.Name, &client.ContactInfo, &client.Address)
		if err != nil {
			if err == sql.ErrNoRows {
				tmpl.Execute(w, "")
				http.Error(w, "Client not found", http.StatusBadRequest)
				return
			}
			tmpl.Execute(w, "")
			http.Error(w, "Failed to get information about client", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, client)
		if err != nil {
			fmt.Errorf("error executing template: %s", err)
		}
	}
}

func ClientSubmitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("Name")
		contactInfo := r.FormValue("ContactInfo")
		address := r.FormValue("Address")

		validName := regexp.MustCompile(`^[a-zA-Z]+ [a-zA-Z]`)
		if len(contactInfo) != 10 || !validName.MatchString(name) {
			http.Redirect(w, r, "submit/404", http.StatusSeeOther)
			http.Error(w, "Unable to add this data", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO client (name, contact_info, address) values (?,?,?)"
		_, err := db.Exec(query, name, contactInfo, address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "submit/success", http.StatusFound)
	}
}
