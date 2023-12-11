package DatabaseHandlers

import (
	"Autoservice/authentication"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC850), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/login" || path == "register" || path == "/login/" || path == "/register" || path == "/register/" {
			next.ServeHTTP(w, r)
			return
		}

		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IsAuthenticated(r *http.Request) bool {
	session, _ := authentication.SessionStore.Get(r, "session-name")
	userID := session.Values["userID"]

	if userID == nil {
		return false
	}
	return true
}
