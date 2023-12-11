package authentication

import (
	"Autoservice/hash"
	"database/sql"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type UserData struct {
	UserID       int    `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserPassword string `json:"user_password"`
}

var SessionStore *sessions.CookieStore

func RegisterHandler(db *sql.DB, hasher *hash.SHA1Hasher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.FormValue("username")
		password := r.FormValue("password")
		if len(login) < 2 {
			http.Error(w, "Too short login", http.StatusBadRequest)
			return
		}
		if len(password) > 14 && len(password) < 5 {
			http.Error(w, "password can't be used", http.StatusBadRequest)
			return
		}

		hashedPassword, err := hasher.Hash(password)
		if err != nil {
			http.Error(w, "Error with password", http.StatusInternalServerError)
			return
		}

		query := "INSERT INTO userinfo (user_login, user_password) VALUES (?, ?)"
		_, err = db.Exec(query, login, hashedPassword)
		if err != nil {
			log.Println(err)
			http.Error(w, "error creating user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func LoginHandler(db *sql.DB, hasher *hash.SHA1Hasher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		hashedPassword, err := hasher.Hash(password)
		if err != nil {
			http.Error(w, "Unable to check password", http.StatusBadRequest)
		}
		var user UserData
		query := "SELECT user_id,user_password FROM userinfo where user_login = ?"
		row := db.QueryRow(query, username)

		err = row.Scan(&user.UserID, &user.UserPassword)
		if err != nil {
			http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
			return
		}
		if hashedPassword != user.UserPassword {
			http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
			return
		}

		session, _ := SessionStore.Get(r, "session-name")
		session.Values["userID"] = user.UserID
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := SessionStore.Get(r, "session-name")
		session.Values = make(map[interface{}]interface{})
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func InitSessionStore() {
	SessionStore = sessions.NewCookieStore([]byte("your-secret-key"))
}
