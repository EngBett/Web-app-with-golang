package main

import (
	"./models"
	"./sessions"
	"./utils"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	models.Init()

	utils.LoadTemplates("templates/*.html")

	r := mux.NewRouter()
	r.HandleFunc("/", AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", AuthRequired(indexPostHandler)).Methods("POST")

	/**
	Login Routes
	*/
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")

	/**
	Registration Routes
	*/
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	//Serve
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

/**
Middleware
*/
func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		_, ok := session.Values["username"]

		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	comments, err := models.GetComments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	templates.ExecuteTemplate(w, "index.html", comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	err := models.PostComment(comment)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/", 302)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	err := models.AuthenticatesUser(username, password)

	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			templates.ExecuteTemplate(w, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			templates.ExecuteTemplate(w, "login.html", "Invalid login")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}

		return
	}

	session, _ := sessions.Store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	err := models.RegisterUser(username, password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	http.Redirect(w, r, "/login", 302)
}
