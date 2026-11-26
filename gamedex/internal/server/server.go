package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/catalog"
)

type Server struct {
	tmpl    *template.Template
	catalog *catalog.Catalog
}

type ListePage struct {
	Games []catalog.Game
	Query string
}

type JeuPage struct {
	Game        catalog.Game
	Tournaments []catalog.Tournament
}

type ErreurPage struct {
	Code    int
	Message string
}

func New(c *catalog.Catalog) *Server {
	return &Server{
		tmpl:    template.Must(template.ParseGlob("web/templates/*.html")),
		catalog: c,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.accueil)
	mux.HandleFunc("GET /jeux", s.liste)
	mux.HandleFunc("GET /jeux/{id}", s.detailJeu)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	return mux
}

func (s *Server) erreur(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	err := s.tmpl.ExecuteTemplate(w, "erreur.html", ErreurPage{Code: code, Message: message})
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) accueil(w http.ResponseWriter, r *http.Request) {
	err := s.tmpl.ExecuteTemplate(w, "accueil.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) liste(w http.ResponseWriter, r *http.Request) {
	page := ListePage{Games: s.catalog.Games(), Query: r.URL.Query().Get("q")}
	err := s.tmpl.ExecuteTemplate(w, "liste.html", page)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		s.erreur(w, http.StatusBadRequest, "Identifiant invalide.")
		return
	}

	game, ok := s.catalog.GameByID(id)
	if !ok {
		s.erreur(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	page := JeuPage{Game: game, Tournaments: s.catalog.TournamentsOf(id)}
	err = s.tmpl.ExecuteTemplate(w, "jeu.html", page)
	if err != nil {
		log.Println(err)
	}
}
