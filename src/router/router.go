package router

import (
	"groupie/controller"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/collection", controller.CollectionHandler)
	mux.HandleFunc("/ressources", controller.RessourcesHandler)
	mux.HandleFunc("/categorie", controller.CategorieHandler)
	mux.HandleFunc("/aPropos", controller.AProposHandler)
	mux.HandleFunc("/ressource", controller.RessourceHandler)
    mux.HandleFunc("/favoris/toggle", controller.ToggleFavoris)


	return mux
}
