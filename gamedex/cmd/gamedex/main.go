package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/server"
)

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	games, err := api.FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux : ", err)
	}
	tournaments, err := api.FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois : ", err)
	}

	cat := catalog.New(games, tournaments)
	srv := server.New(cat)

	fmt.Printf("Serveur sur http://localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), srv.Handler()))
}
