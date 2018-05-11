package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rubinda/biolog"
	log "github.com/sirupsen/logrus"
	//	log "github.com/sirupsen/logrus"
)

// SpeciesHandler je http handler za SpeciesService
type SpeciesHandler struct {
	SpeciesService biolog.SpeciesService
	*chi.Mux
}

// NewSpeciesHandler kreira novega handlerja za vrste in operacije povezane z njimi
func NewSpeciesHandler() *SpeciesHandler {
	sh := &SpeciesHandler{
		Mux: chi.NewRouter(),
	}

	// Prefix do tukaj je ze /api/v1/species, poti pisemo od tega naprej
	// Parametre za iskanje na posameznih endpointih (npr. ?phylum=value) najdemo znotraj
	// Handler funkcij
	sh.Get("/", sh.GetAllSpecies)
	sh.Get("/{gbifKey:[0-9]+}", sh.GetSpeciesByGBIFKey)
	sh.Post("/", sh.CreateSpecies)
	sh.Patch("/{gbifKey:[0-9]+}", sh.UpdateLocalSpecies)
	sh.Delete("/{gbifKey:[0-9]+}", sh.DeleteSpecies)

	sh.Get("/observations", sh.GetObservations)
	sh.Get("/observations/{id:[0-9]+}", sh.GetObservationByID)
	sh.Post("/observations", sh.CreateObservation)
	sh.Patch("/observaions/{id:[0-9]+}", sh.UpdateObservation)
	sh.Delete("/observations/{id:[0-9]+}", sh.DeleteObservation)

	return sh
}

// GetAllSpecies vrne vse vrste, ki so bile popisane in shranjene pri nas
// TODO:
// 	- boljse javljanje napak
func (sh *SpeciesHandler) GetAllSpecies(w http.ResponseWriter, r *http.Request) {
	// Mozne parametre za razne queryje lahko dobimo od tega
	// q := r.URL.Query()
	// family := q.Get("family")

	sps, err := sh.SpeciesService.AllSpecies()
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, sps)

}

// GetSpeciesByGBIFKey vrne podrobnosti o vrsti shranjene pri nas preko kljuca od GBIF
func (sh *SpeciesHandler) GetSpeciesByGBIFKey(w http.ResponseWriter, r *http.Request) {
	gbifKey, parseErr := getIDFromURL(w, r, "gbifKey")
	if parseErr {
		return
	}

	sp, err := sh.SpeciesService.Species(gbifKey)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, sp)
}

// CreateSpecies kreira nov zapis o neki vrsti v naso podatkovno bazo
func (sh *SpeciesHandler) CreateSpecies(w http.ResponseWriter, r *http.Request) {
	var sp biolog.Species

	// Podatki iz telesa
	decErr := json.NewDecoder(r.Body).Decode(&sp)

	// Pri pretvarjanju je prislo do napake
	if decErr != nil {

		switch decErr {
		case io.EOF:
			respondWithError(w, 400, "Telo zahtevka pri kreiranju vrste ne more biti prazno")
		default:
			respondWithError(w, 400, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	// Shrani podatke o novi vrsti
	newSp, err := sh.SpeciesService.CreateSpecies(&sp)

	// Napaka pri kreiranju
	if err != nil {
		log.Error(err)
		respondWithError(w, 400, "Napaka pri ustvarjanju nove vrste")
		return
	}

	// Vrsta uspesno shranjena, vrni novo shranjene podatke
	respondWithJSON(w, 201, newSp)
}

// UpdateLocalSpecies posodobi podatke o lokalno shranjeni vrsti
func (sh *SpeciesHandler) UpdateLocalSpecies(w http.ResponseWriter, r *http.Request) {
	gbifKey, parseErr := getIDFromURL(w, r, "gbifKey")
	if parseErr {
		return
	}

	var sp biolog.Species
	// Podatki iz telesa
	decErr := json.NewDecoder(r.Body).Decode(&sp)

	// Pri pretvarjanju je prislo do napake
	if decErr != nil {

		switch decErr {
		case io.EOF:
			respondWithError(w, 400, "Telo zahtevka pri kreiranju vrste ne more biti prazno")
		default:
			respondWithError(w, 400, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	err := sh.SpeciesService.UpdateSpecies(gbifKey, sp)

	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}

// DeleteSpecies zbrise lokalno shranjeno vrsto
func (sh *SpeciesHandler) DeleteSpecies(w http.ResponseWriter, r *http.Request) {
	gbifKey, parseErr := getIDFromURL(w, r, "gbifKey")
	if parseErr {
		return
	}

	if err := sh.SpeciesService.DeleteSpecies(gbifKey); err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}

// GetObservations vrne vse opazovalne liste
func (sh *SpeciesHandler) GetObservations(w http.ResponseWriter, r *http.Request) {
	obs, err := sh.SpeciesService.Observations()

	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, obs)
}

// GetObservationByID vrne tocno dolocen opazovalni list
func (sh *SpeciesHandler) GetObservationByID(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	ob, err := sh.SpeciesService.Observation(id)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, ob)
}

// CreateObservation ustvari nov Observation za doloceno vrsto
func (sh *SpeciesHandler) CreateObservation(w http.ResponseWriter, r *http.Request) {
	var ob biolog.Observation

	// Podatki iz telesa
	decErr := json.NewDecoder(r.Body).Decode(&ob)

	// Pri pretvarjanju je prislo do napake
	if decErr != nil {

		switch decErr {
		case io.EOF:
			respondWithError(w, 400, "Telo zahtevka pri zapisu opazovanja vrste ne more biti prazno")
		default:
			respondWithError(w, 400, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	// Shrani podatke o novi vrsti
	newOb, err := sh.SpeciesService.CreateObservation(&ob)

	// Napaka pri kreiranju
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// Opazovanje vrste uspesno shranjeno, vrni novo shranjene podatke
	respondWithJSON(w, 201, newOb)
}

// UpdateObservation posodobi dolocen opazovalni list
func (sh *SpeciesHandler) UpdateObservation(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	var ob biolog.Observation
	// Podatki iz telesa
	decErr := json.NewDecoder(r.Body).Decode(&ob)

	// Pri pretvarjanju je prislo do napake
	if decErr != nil {

		switch decErr {
		case io.EOF:
			respondWithError(w, 400, "Telo pri posodabljanju opazovanja ne more biti prazno")
		default:
			respondWithError(w, 400, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	err := sh.SpeciesService.UpdateObservation(id, ob)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}

// DeleteObservation izbrise dolocen opazovalni list
func (sh *SpeciesHandler) DeleteObservation(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	if err := sh.SpeciesService.DeleteObservation(id); err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}
