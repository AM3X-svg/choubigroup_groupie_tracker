package controller

import (
	"encoding/json"
	"fmt"
	"groupie/pages"
	Struct "groupie/struct"
	"io"
	"net/http"
	"time"
)

func renderPage(w http.ResponseWriter, filename string, data any) {
	err := pages.Temp.ExecuteTemplate(w, filename, data)
	if err != nil {
		// Meilleure pratique : logguer l'erreur côté serveur
		fmt.Println("Erreur rendu template :", err)
		http.Error(w, "Erreur rendu template.", http.StatusInternalServerError)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		// Maintenant GetPokedex() retourne bien le slice de Pokémon
		"Pokedex": GetPokedex(),
	}
	renderPage(w, "index.html", data)
}

// GetPokedex retourne un slice de Struct.ApiData, représentant la liste des Pokémon.
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

	// Vérifier le statut de la réponse avant de lire le corps
	if res.StatusCode != http.StatusOK {
		fmt.Println("Erreur statut HTTP :", res.StatusCode)
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erreur lecture corps :", err)
		return nil
	}

	// ⭐️ CORRECTION MAJEURE : On désérialise dans un SLICE de ApiData.
	var pokedex []Struct.ApiData
	if err := json.Unmarshal(body, &pokedex); err != nil {
		fmt.Println("Erreur décodage JSON :", err)
		// Optionnel : Afficher le corps pour aider au débogage du JSON
		// fmt.Println("Corps de la réponse :", string(body))
		return nil
	}

	if len(pokedex) > 0 {
		fmt.Println("✅ Succès ! %d Pokémon décodés. Le premier est : %s\n", len(pokedex), pokedex[0].Name.Fr)
	} else {
		fmt.Println("⚠️ Avertissement : Aucune donnée de Pokémon n'a été décodée (slice vide).")
	}

	// ⭐️ DEUXIÈME CORRECTION : On retourne le slice complet.
	// La ligne `return data.PokedexId` n'est plus nécessaire/pertinente.
	return pokedex
}
