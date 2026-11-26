package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"text/template" 
	"log"
	"net/http"
	"strconv"
	"time"
)

const baseURL = "http://localhost:8081"

var client = &http.Client{Timeout: 3 * time.Second}

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

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
type ListePage struct {
	Games []Game
	Query string
}

type JeuPage struct {
	Game        Game
	Tournaments []Tournament
}

type ErreurPage struct {
	Code    int
	Message string
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

func erreur(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code) 
	err := tmpl.ExecuteTemplate(w, "erreur.html", ErreurPage{Code: code, Message: message})
	if err != nil {
		log.Println(err)
	}
}

func accueil(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "accueil.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func liste(w http.ResponseWriter, r *http.Request) {
	games, err := FetchGames()
	if err != nil {
		log.Println("API :", err)
		erreur(w, http.StatusServiceUnavailable, "L'API est indisponible, réessayez plus tard.")
		return
	}

	q := r.URL.Query().Get("q") 

	err = tmpl.ExecuteTemplate(w, "liste.html", ListePage{Games: games, Query: q})
	if err != nil {
		log.Println(err)
	}
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		erreur(w, http.StatusBadRequest, "Identifiant invalide.")
		return
	}

	games, err := FetchGames()
	if err != nil {
		log.Println("API :", err)
		erreur(w, http.StatusServiceUnavailable, "L'API est indisponible, réessayez plus tard.")
		return
	}

	var found *Game
	for i := range games {
		if games[i].ID == id {
			found = &games[i]
			break
		}
	}
	if found == nil {
		erreur(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	tournaments, err := FetchTournaments()
	if err != nil {
		log.Println("API :", err)
		erreur(w, http.StatusServiceUnavailable, "L'API est indisponible, réessayez plus tard.")
		return
	}

	var ofGame []Tournament
	for _, t := range tournaments {
		if t.GameID == id {
			ofGame = append(ofGame, t)
		}
	}

	err = tmpl.ExecuteTemplate(w, "jeu.html", JeuPage{Game: *found, Tournaments: ofGame})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", liste)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Serveur sur http://localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
