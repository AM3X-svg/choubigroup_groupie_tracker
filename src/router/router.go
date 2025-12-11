package router

import (
	"groupie/controller"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.Home)

	return mux
}
