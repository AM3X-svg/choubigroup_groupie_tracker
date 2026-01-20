package controller

import (
	"encoding/json"
	"fmt"
	"groupie/pages"
	Struct "groupie/struct"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func renderPage(w http.ResponseWriter, filename string, data any) {
	err := pages.Temp.ExecuteTemplate(w, filename, data)
	if err != nil {
		fmt.Println("Erreur rendu template :", err)
		http.Error(w, "Erreur rendu template.", http.StatusInternalServerError)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	pokedex := GetPokedex()
	query := r.FormValue("search")

	var results []Struct.ApiData
	var randomPokemon Struct.ApiData

	if len(pokedex) > 0 {
		randomIndex := rand.Intn(len(pokedex))
		randomPokemon = pokedex[randomIndex]
	}

	if query != "" {
		for _, p := range pokedex {
			if strings.Contains(strings.ToLower(p.Name.Fr), strings.ToLower(query)) {
				results = append(results, p)
			}
		}
	} else {
		results = nil
	}

	data := map[string]interface{}{
		"RandomPokemon": randomPokemon,
		"Pokedex":       results,
		"Query":         query,
	}

	renderPage(w, "index.html", data)
}

func CollectionHandler(w http.ResponseWriter, r *http.Request) {

	allPokedex := GetPokedex()

	query := r.FormValue("search")

	var data []Struct.ApiData

	if query != "" {
		for _, p := range allPokedex {
			if strings.Contains(strings.ToLower(p.Name.Fr), strings.ToLower(query)) {
				data = append(data, p)
			}
		}
	} else {
		data = allPokedex
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/collection.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func RessourcesHandler(w http.ResponseWriter, r *http.Request) {
	data := GetPokedex()
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

func CategorieHandler(w http.ResponseWriter, r *http.Request) {
	data := GetPokedex()
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

func AProposHandler(w http.ResponseWriter, r *http.Request) {
	data := GetPokedex()
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

func GetPokedex() []Struct.ApiData {
	urlApi := "https://tyradex.app/api/v1/pokemon"

	req, err := http.NewRequest(http.MethodGet, urlApi, nil)
	if err != nil {
		fmt.Println("Erreur création requête :", err)
		return nil
	}

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur requête HTTP :", err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Erreur statut HTTP :", res.StatusCode)
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erreur lecture corps :", err)
		return nil
	}

	var pokedex []Struct.ApiData
	if err := json.Unmarshal(body, &pokedex); err != nil {
		fmt.Println("Erreur décodage JSON :", err)
		return nil
	}

	return pokedex
}

func SearchSystem(query string) []Struct.ApiData {
	// 1. L'URL de base (On récupère tout car l'API n'a pas de endpoint /Search)
	urlApi := "https://tyradex.app/api/v1/pokemon"

	req, err := http.NewRequest(http.MethodGet, urlApi, nil)
	if err != nil {
		return nil
	}

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var allPokemon []Struct.ApiData
	json.Unmarshal(body, &allPokemon)

	// 2. Logique de recherche : on filtre manuellement
	var resultats []Struct.ApiData
	for _, p := range allPokemon {
		// On compare le nom français du Pokémon avec la recherche (query)
		if strings.Contains(strings.ToLower(p.Name.Fr), strings.ToLower(query)) {
			resultats = append(resultats, p)
		}
	}

	return resultats
}
