package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8081"

var client = &http.Client{Timeout: 3 * time.Second}

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
}


func main() {
	games, err := FetchGames()
	if err != nil {
		log.Fatal("Impossible de récupérer les jeux : ", err)
	}

	tournaments, err := FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de récupérer les tournois : ", err)
	}

	byGame := make(map[int][]Tournament)
	for _, t := range tournaments {
		byGame[t.GameID] = append(byGame[t.GameID], t)
	}

	for _, g := range games {
		fmt.Printf("%d. %s (%d) - %d tournoi(s)\n", g.ID, g.Name, g.Year, len(byGame[g.ID]))
	}
}

func FetchGames() ([]Game, error) {
	resp, err := client.Get(baseURL + "/api/games")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var games []Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	return games, err
}

func FetchTournaments() ([]Tournament, error) {
	resp, err := client.Get(baseURL + "/api/tournaments")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var tournaments []Tournament
	err = json.NewDecoder(resp.Body).Decode(&tournaments)
	return tournaments, err
}

