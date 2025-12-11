package structure

import (
	"go/types"
)

type Pokemon struct {
	PokedexID   int          `json:"pokedex_id"`
	Generation  int          `json:"generation"`
	Category    string       `json:"category"`
	Name        string       `json:"name"`
	Sprite      string       `json:"sprite"`
	Types       []types.Info `json:"types"`
	Talents     []Talent     `json:"talents"`
	Stats       Stats        `json:"stats"`
	Resistances []Resistance `json:"resistances"`
	Evolution   Evolution    `json:"evolution"`
	Weight      string       `json:"weight"`
	Height      string       `json:"height"`
	Egggroup    []string     `json:"egg_group"`
	Sexe        Sexe         `json:"sexe"`
	Catchrate   int          `json:"catch_rate"`
	Level100    int          `json:"level_100"`
	Formes      []string     `json:"formes"`
}

type LacalizedName struct {
	Fr string `json:"fr"`
	En string `json:"en"`
	Jp string `json:"jp"`
}

type Sprites struct {
	Regular string  `json:"regular"`
	Shiny   string  `json:"shiny"`
	Gmax    *string `json:"gmax"`
}

type Typeinfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Talent struct {
	Name string `json:"name"`
	TC   bool   `json:"TC"`
}

type Stats struct {
	HP        int `json:"HP"`
	Attack    int `json:"Attack"`
	Defense   int `json:"Defense"`
	SpAttack  int `json:"SpAttack"`
	SpDefense int `json:"SpDefense"`
	Speed     int `json:"Speed"`
}

type Resistance struct {
	Name       string  `json:"name"`
	Multiplier float64 `json:"multiplier"`
}

type Evolution struct {
	Pre  []Evolutionlink `json:"pre"`
	Next []Evolutionlink `json:"next"`
	Mega []Evolutionlink `json:"mega"`
}

type Evolutionlink struct {
	PokedexID int    `json:"pokedex_id"`
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type Sexe struct {
	Male   int `json:"male"`
	Female int `json:"female"`
}
