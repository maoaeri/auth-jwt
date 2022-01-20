package api

import (
	"html/template"
	"log"
	"myapp/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func CreateRouter() {
	router = mux.NewRouter()
}

func InitRouter() {
	router.HandleFunc("/register", Register)                    //.Methods("POST")
	router.HandleFunc("/login", LogIn)                          //.Methods("POST")
	router.HandleFunc("/index", middleware.IsAuthorized(Index)) //.Methods("GET")
	router.HandleFunc("/logout", LogOut)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../pkg/template/Home.html"))
		tmpl.Execute(w, true)
	})
}

func Run() {
	CreateRouter()
	InitRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
