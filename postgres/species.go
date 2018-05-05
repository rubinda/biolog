package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // dodatek PostgreSQL
	"github.com/rubinda/biolog"
)

// SpeciesService predstavlja PostgreSQL implementacijo od biolog.SpeciesService
type SpeciesService struct {
	DB *sqlx.DB
}

// Species vrne doloceno vrsto, ki je shranjena pri nas
func (s *SpeciesService) Species(id int) (*biolog.Species, error) {
	return nil, errors.New("Not implemented (yet)")
}

// CreateSpecies shrani podatke o doloceni vrsti v naso bazo in vrne dodeljen id
func (s *SpeciesService) CreateSpecies(sp *biolog.Species) (int64, error) {
	return -1, errors.New("Not implemented")
}

// Observation vrne zapis z dolocenim ID
func (s *SpeciesService) Observation(id int) (*biolog.Observation, error) {
	return nil, errors.New("Not implemented")
}

// Observations vrne vse podane zapise o opazenih vrstah
func (s *SpeciesService) Observations() ([]*biolog.Observation, error) {
	return nil, errors.New("Not implemented")
}

// CreateObservation kreira nov zapis o opazeni vrsti
func (s *SpeciesService) CreateObservation(o *biolog.Observation) (int, error) {
	return -1, errors.New("Not implemented")
}

// DeleteObservation zbrise dolocen zapis o opazeni vrsti
func (s *SpeciesService) DeleteObservation(id int) error {
	return errors.New("Not implemented")
}

// UpdateObservation posodobi opazovalni list, ki ima enak ID
// Nove podatke preberemo iz slovarja, pri cemer so kljuci enaki imenom atributov
func (s *SpeciesService) UpdateObservation(o map[string]string) error {
	return errors.New("Not implemented")
}

// ConservationStatus vrne podatke o dolocenem statusu ogrozenosti
func (s *SpeciesService) ConservationStatus(id int) (*biolog.ConservationStatus, error) {
	return nil, errors.New("Not implemented")
}

// ConservationStatuses vrne vse mozne statuse ogrozenosti za doloceno vrsto
func (s *SpeciesService) ConservationStatuses() ([]*biolog.ConservationStatus, error) {
	return nil, errors.New("Not implemented")
}
