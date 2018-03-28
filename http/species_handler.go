package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rubinda/biolog"
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
func (sh *SpeciesHandler) GetAllSpecies(w http.ResponseWriter, r *http.Request) {
	// Mozne parametre za razne queryje lahko dobimo od tega
	q := r.URL.Query()

	family := q.Get("family")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented %s", family)
	// TODO implement the method
}

// GetSpeciesByGBIFKey vrne podrobnosti o vrsti shranjene pri nas preko kljuca od GBIF
func (sh *SpeciesHandler) GetSpeciesByGBIFKey(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// CreateSpecies kreira nov zapis o neki vrsti v naso podatkovno bazo
func (sh *SpeciesHandler) CreateSpecies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// UpdateLocalSpecies posodobi podatke o lokalno shranjeni vrsti
func (sh *SpeciesHandler) UpdateLocalSpecies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// DeleteSpecies zbrise lokalno shranjeno vrsto
func (sh *SpeciesHandler) DeleteSpecies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetObservations vrne vse opazovalne liste
func (sh *SpeciesHandler) GetObservations(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetObservationByID vrne tocno dolocen opazovalni list
func (sh *SpeciesHandler) GetObservationByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// CreateObservation ustvari nov Observation za doloceno vrsto
func (sh *SpeciesHandler) CreateObservation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// UpdateObservation posodobi dolocen opazovalni list
func (sh *SpeciesHandler) UpdateObservation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// DeleteObservation izbrise dolocen opazovalni list
func (sh *SpeciesHandler) DeleteObservation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}
