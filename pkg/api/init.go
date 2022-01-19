package api

import (
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
}

func Run() {
	CreateRouter()
	InitRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
