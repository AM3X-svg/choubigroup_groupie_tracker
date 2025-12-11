package main

import (
	"fmt"
	initTemp "groupie/pages"
	"groupie/router"
	"net/http"
)

func main() {
	initTemp.Init()

	r := router.New()

	// Sert les fichiers statiques (CSS, JS, images, etc.)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	// Couleurs ANSI
	blue := "\033[34m"
	reset := "\033[0m"

	fmt.Println("Serveur démarré sur", blue+"http://localhost:8080"+reset)

	http.ListenAndServe(":8080", r)
}
