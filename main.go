package main

import (
	"Autoservice/DatabaseHandlers"
	"Autoservice/authentication"
	"Autoservice/database"
	"Autoservice/hash"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := database.MySqlConnect(database.ConnectionInfo{
		User:   "root",
		Passwd: "1234",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "mydb",
	})
	if err != nil {
		log.Fatal(err)
	}
	hasher := hash.NewSHA1Hasher("salt-salt")
	router := mux.NewRouter()
	authentication.InitSessionStore()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	router.Use(DatabaseHandlers.LoggingMiddleware)
	router.Use(DatabaseHandlers.AuthenticationMiddleware)
	router.HandleFunc("/404", DatabaseHandlers.ErrorPageHandler)
	router.HandleFunc("/", DatabaseHandlers.HomePageHandler)
	login := router.PathPrefix("/login").Subrouter()
	login.HandleFunc("", DatabaseHandlers.LoginPageHandler)
	login.HandleFunc("/", authentication.LoginHandler(db, hasher))

	register := router.PathPrefix("/register").Subrouter()
	register.HandleFunc("", DatabaseHandlers.RegisterPageHandler)
	register.HandleFunc("/", authentication.RegisterHandler(db, hasher))

	auto := router.PathPrefix("/auto").Subrouter()
	auto.Use(DatabaseHandlers.AuthenticationMiddleware)
	auto.HandleFunc("", DatabaseHandlers.AutoPageHandler)
	auto.HandleFunc("/", DatabaseHandlers.AutoHandlerHtml(db))

	client := router.PathPrefix("/client").Subrouter()
	client.Use(DatabaseHandlers.AuthenticationMiddleware)
	client.HandleFunc("/submit/404", DatabaseHandlers.ErrorPageHandler)
	client.HandleFunc("", DatabaseHandlers.ClientPageHandler)
	client.HandleFunc("/", DatabaseHandlers.ClientHandlerHtml(db))
	client.HandleFunc("/submit", DatabaseHandlers.ClientSubmitHandler(db))
	client.HandleFunc("/submit/success", DatabaseHandlers.SuccessPageHandler)

	order := router.PathPrefix("/order").Subrouter()
	order.Use(DatabaseHandlers.AuthenticationMiddleware)
	order.HandleFunc("", DatabaseHandlers.OrderPageHandler)
	order.HandleFunc("/", DatabaseHandlers.GetInfoOrder(db))

	exit := router.PathPrefix("/logout").Subrouter()
	exit.HandleFunc("", authentication.LogoutHandler())

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Errorf("Error listening: %s", err)
	}
	defer db.Close()

}
