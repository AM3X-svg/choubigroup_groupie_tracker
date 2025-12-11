package structure

type ApiData struct {
	PokedexId   int           `json:"pokedex_id"`
	Generation  int           `json:"generation"`
	Categorie   string        `json:"categorie"`
	Name        Name          `json:"name"`
	Sprites     Sprites       `json:"sprites"`
	Types       []Types       `json:"types"`
	Talents     []Talents     `json:"talents"`
	Stats       Stats         `json:"stats"`
	Resistances []Resistances `json:"resistances"`
	Evolution   Evolution     `json:"evolution"`
	Height      string        `json:"height"`
	Weight      string        `json:"weight"`
	EggGroups   []string      `json:"egg_groups"`
	Sexe        Sexe          `json:"sexe"`
	CatchRate   int           `json:"catch_rate"`
	Level100    int           `json:"level_100"`
	Formes      []Formes      `json:"formes"`
}

type Name struct {
	Fr string `json:"fr"`
	En string `json:"en"`
	Jp string `json:"jp"`
}

type Sprites struct {
	Regular string `json:"regular"`
	Shiny   string `json:"shiny"`
	Gmax    Gmax   `json:"gmax"`
}

type Gmax struct {
	Regular string `json:"regular"`
	Shiny   string `json:"shiny"`
}

type Types struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Talents struct {
	Name string `json:"name"`
	Tc   bool   `json:"tc"`
}

type Stats struct {
	Hp     int `json:"hp"`
	Atk    int `json:"atk"`
	Def    int `json:"def"`
	SpeAtk int `json:"spe_atk"`
	SpeDef int `json:"spe_def"`
	Vit    int `json:"vit"`
}

type Resistances struct {
	Name       string  `json:"name"`
	Multiplier float64 `json:"multiplier"`
}

type Evolution struct {
	Pre  []Pre  `json:"pre"`
	Next []Next `json:"next"`
	Mega []Mega `json:"mega"`
}

type Pre struct {
	PokedexId int    `json:"pokedex_id"`
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type Next struct {
	PokedexId int    `json:"pokedex_id"`
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type Mega struct {
	Orbe        string      `json:"orbe"`
	SpritesMega SpritesMega `json:"sprites"`
}

type SpritesMega struct {
	Regular string `json:"regular"`
	Shiny   string `json:"shiny"`
}

type Sexe struct {
	Male   float64 `json:"male"`
	Female float64 `json:"female"`
}

type Formes struct {
	Region string `json:"region"`
	Name   Name   `json:"name"`
}
