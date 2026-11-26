package catalog

type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

type Tournament struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	GameID int    `json:"game_id"`
	City   string `json:"city"`
	Date   string `json:"date"`
}

type Catalog struct {
	games  []Game
	byGame map[int][]Tournament
}

func New(games []Game, tournaments []Tournament) *Catalog {
	byGame := make(map[int][]Tournament)
	for _, t := range tournaments {
		byGame[t.GameID] = append(byGame[t.GameID], t)
	}
	return &Catalog{games: games, byGame: byGame}
}

func (c *Catalog) Games() []Game {
	return c.games
}

func (c *Catalog) GameByID(id int) (Game, bool) {
	for _, g := range c.games {
		if g.ID == id {
			return g, true
		}
	}
	return Game{}, false
}

func (c *Catalog) TournamentsOf(gameID int) []Tournament {
	return c.byGame[gameID]
}
