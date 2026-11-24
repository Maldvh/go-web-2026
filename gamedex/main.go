package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"log"
)



func accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>Gamedex</title>
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <h1>Gamedex</h1>
  <p><a href="/jeux">Voir la liste des jeux</a></p>
</body>
</html>`)
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux")
}
func main() {
	port := flag.String("port", "8080", "port du serveur")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	addr := ":" + *port
	fmt.Printf("Serveur sur http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr,mux))
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	if id > 12 {
		http.NotFound(w, r) 
		return
	}
	fmt.Fprintf(w, "jeu %d\n", id)
}

