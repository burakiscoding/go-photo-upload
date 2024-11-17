package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type HomePageData struct {
	Images []string
}

type Handler struct {
	tmpl *template.Template
}

func NewHandler(tmpl *template.Template) *Handler {
	return &Handler{tmpl: tmpl}
}

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)
	if err := r.ParseMultipartForm(15 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join("uploads", header.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "image", header.Filename)
}

func (h *Handler) HandleGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "bad image name", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, filepath.Join("uploads", name))
}

func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	photo := HomePageData{}

	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		photo.Images = append(photo.Images, file.Name())
	}

	h.tmpl.Execute(w, photo)
}
