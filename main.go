package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	handler := NewHandler(tmpl)
	router.HandleFunc("/", handler.HandleHomePage).Methods(http.MethodGet)
	router.HandleFunc("/upload", handler.HandleUpload).Methods(http.MethodPost)
	router.HandleFunc("/images/{name}", handler.HandleGetImage).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}
