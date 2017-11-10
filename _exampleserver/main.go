package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret"))

func ensureLogin(r *http.Request) (string, bool) {
	session, err := store.Get(r, "user")
	if err != nil {
		return "", false
	}
	nameintf, ok := session.Values["name"]
	if !ok {
		return "", false
	}
	name, ok := nameintf.(string)
	return name, ok
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	session, err := store.New(r, "user")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session.Values["name"] = name
	if err := store.Save(r, w, session); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<body>
				<form method="post" action="/login">
					<input type="text" name="name">
					<input type="submit" value="Login">
				</form>
			</body>
		</html>
	`))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLogin(w, r)
	case http.MethodPost:
		postLogin(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	delete(session.Values, "name")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	name, ok := ensureLogin(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	fmt.Fprintf(w, `
		<html>
			<body>
				<h2> Welcome %s!</h2>
				<a href="/logout">Logout</a>
			</body>
		</html>
	`, name)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/", getIndex)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
