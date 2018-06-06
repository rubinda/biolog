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

// SpeciesGbifKey model
//
// Za iskanje po lokalno shranjenih vrstah
// swagger:parameters getSpeciesbyGbifKey deleteSpecies updateSpecies
type SpeciesGbifKey struct {
	// in: path
	// required: true
	GbifKey int `json:"gbifKey"`
}

// SpeciesBodyParams model
//
// Pri virih, ki v telesu zahtevajo podatke o vrsti
// swagger:parameters createSpecies
type SpeciesBodyParams struct {
	// in: body
	// required: true
	Payload *biolog.Species `json:"species"`
}

// ObservationID model.
//
// Se uporablja za vire, ki se navezeujejo na opazanja preko IDjev
// swagger:parameters getObservationByID deleteObservation updateObservation
type ObservationID struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// ObservationBodyParams model.
//
// Se uporablja pri virih, ki pricakujejo podatke o opazeni vrsti v telesu zahtevka
// swagger:parameters createObservation
type ObservationBodyParams struct {
	// in: body
	// required: true
	Paylod *biolog.Observation `json:"observation"`
}

// NewSpeciesHandler kreira novega handlerja za vrste in operacije povezane z njimi
func NewSpeciesHandler() *SpeciesHandler {
	sh := &SpeciesHandler{
		Mux: chi.NewRouter(),
	}

	// Prefix do tukaj je ze /api/v1/species, poti pisemo od tega naprej
	// Parametre za iskanje na posameznih endpointih (npr. ?phylum=value) najdemo znotraj
	// Handler funkcij

	// swagger:route GET /species species getSpecies
	//
	// Pridobi vse lokalno shranjene vrste
	//
	// Responses:
	//		200: []species
	sh.Get("/", sh.GetAllSpecies)

	// swagger:route POST /species species createSpecies
	//
	// Ustvari nov zapis o podatkah neke vrste
	//
	// Responses:
	// 		201: species
	sh.Post("/", sh.CreateSpecies)

	// Zdruzi vse podoperacije, ki zahtevajo GBIF Key v URL
	// TODO:
	//	- pridobi ID iz URL preko middleware
	sh.Route("/{gbifKey:[0-9]+}", func(r chi.Router) {
		// swagger:route GET /species/{gbifKey} species getSpeciesByGbifKey
		//
		// Pridobi podrobnosti o vrsti preko GBIF kljuca
		//
		// Responses:
		//		200: species
		r.Get("/", sh.GetSpeciesByGBIFKey)

		// swagger:route PATCH /species/{gbifKey} species updateSpecies
		//
		// Posodobi podakte o shranjeni vrsti
		//
		// Responses:
		//		204:
		r.Patch("/", sh.UpdateLocalSpecies)

		// swagger:route DELETE /species/{gbifKey} species deleteSpecies
		//
		// Zbrise shranjeno vrsto
		//
		// Responses:
		//		204:
		r.Delete("/", sh.DeleteSpecies)
	})

	// Podpoti na /observations
	sh.Route("/observations", func(r chi.Router) {
		// swagger:route GET /species/observations observations getObservations
		//
		// Pridobi vsa javna opazanja
		//
		// Responses:
		//		200: []observation
		r.Get("/", sh.GetObservations)

		// swagger:route POST /species/observations observations createObservation
		//
		// Ustvari nov zapis o opazeni vrsti
		//
		// Responses:
		//		201: observation
		r.Post("/", sh.CreateObservation)

		// TODO:
		//	- pridobi ID iz URL preko middleware (r.Use(GetObservationIDCtx) {...})
		r.Route("/{id:[0-9]+}", func(r chi.Router) {
			// swagger:route GET /species/observations/{id} observations getObservationsByID
			//
			// Pridobi opazanje s podanim IDjem
			//
			// Responses:
			// 		200: observation
			r.Get("/", sh.GetObservationByID)

			// swagger:route PATCH /species/observations/{id} observations updateObservation
			//
			// Posodobi podatke o opazovalnem listu
			//
			// Responses:
			//		204:
			r.Patch("/", sh.UpdateObservation)

			// swagger:route DELETE /species/observations/{id} observations deleteObservation
			//
			// Zbrise podatek o opazeni vrsti
			//
			// Responses:
			// 		204:
			r.Delete("/", sh.DeleteObservation)
		})

	})

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
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sps)

}

// GetSpeciesByGBIFKey vrne podrobnosti o vrsti shranjene pri nas preko kljuca od GBIF
func (sh *SpeciesHandler) GetSpeciesByGBIFKey(w http.ResponseWriter, r *http.Request) {
	gbifKey, parseErr := getIDFromURL(w, r, "gbifKey")
	if parseErr {
		return
	}

	sp, err := sh.SpeciesService.Species(gbifKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sp)
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
			respondWithError(w, http.StatusBadRequest, "Telo zahtevka pri kreiranju vrste ne more biti prazno")
		default:
			respondWithError(w, http.StatusBadRequest, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	// Shrani podatke o novi vrsti
	newSp, err := sh.SpeciesService.CreateSpecies(&sp)

	// Napaka pri kreiranju
	if err != nil {
		log.Error(err)
		respondWithError(w, http.StatusBadRequest, "Napaka pri ustvarjanju nove vrste")
		return
	}

	// Vrsta uspesno shranjena, vrni novo shranjene podatke
	respondWithJSON(w, http.StatusCreated, newSp)
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
			respondWithError(w, http.StatusBadRequest, "Telo zahtevka pri kreiranju vrste ne more biti prazno")
		default:
			respondWithError(w, http.StatusBadRequest, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	err := sh.SpeciesService.UpdateSpecies(gbifKey, sp)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// DeleteSpecies zbrise lokalno shranjeno vrsto
func (sh *SpeciesHandler) DeleteSpecies(w http.ResponseWriter, r *http.Request) {
	gbifKey, parseErr := getIDFromURL(w, r, "gbifKey")
	if parseErr {
		return
	}

	if err := sh.SpeciesService.DeleteSpecies(gbifKey); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// GetObservations vrne vse opazovalne liste
func (sh *SpeciesHandler) GetObservations(w http.ResponseWriter, r *http.Request) {
	obs, err := sh.SpeciesService.Observations()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, obs)
}

// GetObservationByID vrne tocno dolocen opazovalni list
func (sh *SpeciesHandler) GetObservationByID(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	ob, err := sh.SpeciesService.Observation(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ob)
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
			respondWithError(w, http.StatusBadRequest, "Telo zahtevka pri zapisu opazovanja vrste ne more biti prazno")
		default:
			respondWithError(w, http.StatusBadRequest, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	// Shrani podatke o novi vrsti
	newOb, err := sh.SpeciesService.CreateObservation(&ob)

	// Napaka pri kreiranju
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Opazovanje vrste uspesno shranjeno, vrni novo shranjene podatke
	respondWithJSON(w, http.StatusCreated, newOb)
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
			respondWithError(w, http.StatusBadRequest, "Telo pri posodabljanju opazovanja ne more biti prazno")
		default:
			respondWithError(w, http.StatusBadRequest, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}

	err := sh.SpeciesService.UpdateObservation(id, ob)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// DeleteObservation izbrise dolocen opazovalni list
func (sh *SpeciesHandler) DeleteObservation(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	if err := sh.SpeciesService.DeleteObservation(id); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
