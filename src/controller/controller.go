package controller

import (
	"encoding/json"
	"fmt"
	"groupie/pages"
	Struct "groupie/struct"
	"io"
	"math/rand"
	"net/http"
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

	var randomPokemon Struct.ApiData

	if len(pokedex) > 0 {
		randomIndex := rand.Intn(len(pokedex))
		randomPokemon = pokedex[randomIndex]
	} else {
		http.Error(w, "Impossible de charger les données du Pokédex.", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"RandomPokemon": randomPokemon,
		"Pokedex":       pokedex,
	}
	renderPage(w, "index.html", data)
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
