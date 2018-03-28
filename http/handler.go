package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rubinda/biolog"
)

// Handler je kolekcija vseh nasih service handler
type Handler struct {
	UserHandler    *UserHandler
	SpeciesHandler *SpeciesHandler
	*chi.Mux
}

// NewRootHandler ustvari starsa vseh ostalih handlerjev, nosi tudi primarni Router
func NewRootHandler(us biolog.UserService, ss biolog.SpeciesService) *Handler {
	h := &Handler{
		Mux: chi.NewRouter(),
	}

	// A good base middleware stack
	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	// Timeout na zahteve
	h.Use(middleware.Timeout(60 * time.Second))
	// Nastavimo predpono za api
	h.Route("/api/v1", func(r chi.Router) {

		// Podpoti za endpoint '/users'
		h.UserHandler = NewUserHandler()
		h.UserHandler.UserService = us
		r.Mount("/users", h.UserHandler)

		// Podpoti za endpoint '/species'
		h.SpeciesHandler = NewSpeciesHandler()
		h.SpeciesHandler.SpeciesService = ss
		r.Mount("/species", h.SpeciesHandler)
	})

	return h
}

// RespondWithError vrne napako kot odgovor na http request s podanimi podrobnostmi
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON vrne JSON kot odgovor na zahtevo. Parametra sta http koda odgovora in telo
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
