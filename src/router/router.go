package router

import (
	"groupie/controller"
	"html/template"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/collection", collectionHandler)
	mux.HandleFunc("/ressources", ressourcesHandler)
	mux.HandleFunc("/categorie", categorieHandler)
	mux.HandleFunc("/aPropos", aProposHandler)

	return mux
}

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.GetPokedex()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/collection.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ressourcesHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.GetPokedex()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/ressources.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func categorieHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.GetPokedex()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/categorie.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func aProposHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.GetPokedex()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/aPropos.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
