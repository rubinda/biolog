package http

import (
	"encoding/json"
	"net/http"
	"strconv"
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
// FIXME:
// 	- moznost dodajanja lastnih headerjev
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetIDFromURL pridobi veljavno 32 bitno stevilo, pri cemer podamo ime parametra, ce pride do napake obvesti odjemalca
// Vraca veljaven id (stevilo int32) in vrednost ali je prislo do napake
func getIDFromURL(w http.ResponseWriter, r *http.Request, parameter string) (int, bool) {
	// Pridobi ID iz URL in ga pretvori v stevilo (Router poskrbi da je na tej poti vedno stevilka,
	// zato lahko napako ignoriramo)
	id64, parErr := strconv.ParseInt(chi.URLParam(r, parameter), 10, 32)
	id := int(id64)

	// Pri pretvarjanju v Integer (od 0 do 2^31 -1) je prislo do napake (najverjetneje je overflow)
	if parErr != nil {
		e := parErr.(*strconv.NumError)
		// Obvesti, da ID ni v veljavnem obsegu
		if e.Err == strconv.ErrRange {
			respondWithError(w, 400, "Neveljaven ID: izven obsega")
			// Prislo je do druge napake
		} else {
			respondWithError(w, 400, parErr.Error())
		}
		return 0, true
	}

	return id, false
}
