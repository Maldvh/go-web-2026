package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gamedex/internal/catalog"
)

const baseURL = "http://localhost:8081"

var client = &http.Client{Timeout: 3 * time.Second}

func FetchGames() ([]catalog.Game, error) {
	resp, err := client.Get(baseURL + "/api/games")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var games []catalog.Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	return games, err
}

func FetchTournaments() ([]catalog.Tournament, error) {
	resp, err := client.Get(baseURL + "/api/tournaments")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var tournaments []catalog.Tournament
	err = json.NewDecoder(resp.Body).Decode(&tournaments)
	return tournaments, err
}
