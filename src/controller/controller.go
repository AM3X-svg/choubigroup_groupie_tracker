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
	"net/url"
    "strconv"

)

// PageData contient tout ce dont tes templates ont besoin
type PageData struct {
	Pokedex       []Struct.ApiData
	RandomPokemon Struct.ApiData
	Query         string
	CollectionJS  template.JS
	Favorites     []Struct.ApiData // Ajouté pour corriger l'erreur du template index.html

	Pokemon    Struct.ApiData
    IsFavorite bool

}

// renderPage gère l'affichage des fichiers HTML
func renderPage(w http.ResponseWriter, filename string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := pages.Temp.ExecuteTemplate(w, filename, data)
	if err != nil {
		fmt.Printf("Erreur rendu template %s : %v\n", filename, err)
		// On ne fait pas http.Error ici car le Header est déjà envoyé
		return
	}
}

// GetPokedex récupère les données brutes de l'API Tyradex
func GetPokedex() []Struct.ApiData {
	urlApi := "https://tyradex.app/api/v1/pokemon"
	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Get(urlApi)
	if err != nil {
		fmt.Println("Erreur API :", err)
		return nil
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var pokedex []Struct.ApiData
	json.Unmarshal(body, &pokedex)

	return pokedex
}

// Home : Page d'accueil (index.html)
func Home(w http.ResponseWriter, r *http.Request) {
	pokedex := GetPokedex()
	query := strings.ToLower(r.FormValue("search")) // Récupère la valeur du champ 'search'

	var results []Struct.ApiData
	var randomPokemon Struct.ApiData

	// 1. Logique du Pokémon Aléatoire
	if len(pokedex) > 0 {
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		randomPokemon = pokedex[rng.Intn(len(pokedex)-1)+1]
	}

	// 2. Logique de Recherche
	if query != "" {
		for _, p := range pokedex {
			// On cherche dans le nom français
			if strings.Contains(strings.ToLower(p.Name.Fr), query) {
				results = append(results, p)
			}
		}
	}

	// 3. On envoie tout au template
	data := PageData{
		RandomPokemon: randomPokemon,
		Pokedex:       results, // Les résultats de recherche
		Query:         query,   // On renvoie la query pour l'afficher dans l'input
		Favorites:     []Struct.ApiData{},
	}

	renderPage(w, "index.html", data)
}

// CollectionHandler : Page de la collection avec JS de filtrage
func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	allPokedex := GetPokedex()
	if len(allPokedex) == 0 {
		http.Error(w, "API indisponible", http.StatusServiceUnavailable)
		return
	}

	// EXCLUSION DU PREMIER POKEMON (Index 0)
	// On ne garde que les pokémons du 2ème au dernier
	pokedex := allPokedex[1:]

	query := strings.ToLower(r.FormValue("search"))
	var results []Struct.ApiData

	if query != "" {
		for _, p := range pokedex {
			if strings.Contains(strings.ToLower(p.Name.Fr), query) {
				results = append(results, p)
			}
		}
	} else {
		results = pokedex
	}

	data := PageData{
		Pokedex:      results,
		CollectionJS: GetCollectionJS(),
		Favorites:    []Struct.ApiData{},
	}

	renderPage(w, "collection.html", data)
}


// CategorieHandler : Page des catégories
func CategorieHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "categorie.html", PageData{Pokedex: GetPokedex()})
}

// AProposHandler : Page à propos
func AProposHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "aPropos.html", PageData{})
}

// GetCollectionJS : Retourne le JavaScript de filtrage dynamique
func GetCollectionJS() template.JS {
	return template.JS(`
    const filterCheckboxes = document.querySelectorAll('.filter-checkbox');
    const resetBtn = document.getElementById('resetFilters'); // Assure-toi que l'ID est identique dans le HTML

    function applyFilters() {
        const selectedTypes = Array.from(document.querySelectorAll('.filter-energie .filter-checkbox:checked'))
            .map(cb => cb.value);
        const selectedGens = Array.from(document.querySelectorAll('.filter-gen .filter-checkbox:checked'))
            .map(cb => parseInt(cb.value));

        const cards = document.querySelectorAll('.pokemon-sprite');

        cards.forEach(card => {
            const cardTypes = card.getAttribute('data-types') ? card.getAttribute('data-types').split(',') : [];
            const cardGen = parseInt(card.getAttribute('data-gen'));

            // Logique ET pour les types, OU pour les générations
            const typeMatch = selectedTypes.length === 0 || selectedTypes.every(t => cardTypes.includes(t));
            const genMatch = selectedGens.length === 0 || selectedGens.includes(cardGen);

            card.style.display = (typeMatch && genMatch) ? 'block' : 'none';
        });
    }

    // Fonction de réinitialisation
    function resetAll() {
        // 1. Décocher toutes les cases
        filterCheckboxes.forEach(cb => cb.checked = false);

        // 2. Refermer les menus détails (optionnel mais plus propre)
        document.querySelectorAll('details').forEach(det => det.removeAttribute('open'));

        // 3. Relancer le filtrage (selectedTypes sera vide, donc tout s'affichera)
        applyFilters();
    }

    // Écouteurs d'événements
    filterCheckboxes.forEach(cb => cb.addEventListener('change', applyFilters));

    if (resetBtn) {
        resetBtn.addEventListener('click', (e) => {
            e.preventDefault(); // Empêche un éventuel rechargement de page
            resetAll();
        });
    }
    `)
}

const favoritesCookieName = "favorites"

// lit le cookie et renvoie un set d'IDs favoris
func readFavorites(r *http.Request) map[int]bool {
	out := map[int]bool{}

	c, err := r.Cookie(favoritesCookieName)
	if err != nil || c.Value == "" {
		return out
	}

	raw, err := url.QueryUnescape(c.Value)
	if err != nil {
		return out
	}

	var tmp map[string]bool
	if err := json.Unmarshal([]byte(raw), &tmp); err != nil {
		return out
	}

	for k, v := range tmp {
		id, err := strconv.Atoi(k)
		if err == nil {
			out[id] = v
		}
	}
	return out
}

func writeFavorites(w http.ResponseWriter, favs map[int]bool) {
	tmp := map[string]bool{}
	for k, v := range favs {
		tmp[strconv.Itoa(k)] = v
	}

	b, _ := json.Marshal(tmp)
	escaped := url.QueryEscape(string(b))

	http.SetCookie(w, &http.Cookie{
		Name:     favoritesCookieName,
		Value:    escaped,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// transforme le set en slice ApiData (utile si tu veux une page Collection = favoris)
func favoritesToSlice(pokedex []Struct.ApiData, favs map[int]bool) []Struct.ApiData {
	out := make([]Struct.ApiData, 0)
	for _, p := range pokedex {
		if favs[p.PokedexId] {
			out = append(out, p)
		}
	}
	return out
}

// RessourceHandler : Page détail d’un Pokémon (ressource.html)
func RessourceHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	pokedex := GetPokedex()
	if len(pokedex) == 0 {
		http.Error(w, "API indisponible", http.StatusServiceUnavailable)
		return
	}

	// trouver le pokémon par PokedexId
	var found Struct.ApiData
	ok := false
	for _, p := range pokedex {
		if p.PokedexId == id {
			found = p
			ok = true
			break
		}
	}

	if !ok {
		http.Error(w, "pokemon not found", http.StatusNotFound)
		return
	}

	favs := readFavorites(r)

	data := PageData{
		Pokedex:     pokedex,           // pas obligatoire, mais utile si tu veux l'utiliser ailleurs
		Pokemon:     found,             // IMPORTANT : utilisé dans ressource.html
		IsFavorite:  favs[found.PokedexId],
		Favorites:   favoritesToSlice(pokedex, favs),
	}

	renderPage(w, "ressource.html", data)
}

// ToggleFavoris : Ajoute/retire un pokémon des favoris, puis redirect
func ToggleFavoris(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	idStr := r.Form.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	favs := readFavorites(r)

	// toggle
	if favs[id] {
		delete(favs, id)
	} else {
		favs[id] = true
	}

	writeFavorites(w, favs)

	// retour à la page détail
	http.Redirect(w, r, "/ressource?id="+strconv.Itoa(id), http.StatusSeeOther)
}

// RessourcesHandler : Page des ressources (liste + recherche)
func RessourcesHandler(w http.ResponseWriter, r *http.Request) {
	all := GetPokedex()
	if len(all) == 0 {
		http.Error(w, "API indisponible", http.StatusServiceUnavailable)
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	results := all

	if query != "" {
		tmp := make([]Struct.ApiData, 0)
		for _, p := range all {
			if strings.Contains(strings.ToLower(p.Name.Fr), query) {
				tmp = append(tmp, p)
			}
		}
		results = tmp
	}

	favs := readFavorites(r)

	data := PageData{
		Pokedex:   results,
		Query:     query,
		Favorites: favoritesToSlice(all, favs),
	}

	renderPage(w, "ressources.html", data)
}
