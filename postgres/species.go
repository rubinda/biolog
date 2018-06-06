package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // dodatek PostgreSQL
	"github.com/rubinda/biolog"
	//	log "github.com/sirupsen/logrus"
)

// SpeciesService predstavlja PostgreSQL implementacijo od biolog.SpeciesService
type SpeciesService struct {
	DB *sqlx.DB
}

// Species vrne doloceno vrsto, ki je shranjena pri nas, sklicujemo se na GBIF id
func (s *SpeciesService) Species(id int) (*biolog.Species, error) {
	stmt := `SELECT * FROM species WHERE id = $1 LIMIT 1`
	spec := &biolog.Species{}
	if getErr := s.DB.Get(spec, stmt, id); getErr != nil {
		if getErr == sql.ErrNoRows {
			return nil, errors.New("Vrsta s tem GBIF ID ne obstaja")
		}
		return nil, getErr
	}
	return spec, nil
}

// AllSpecies vrne vse vrste, ki so shranjene pri nas
func (s *SpeciesService) AllSpecies() ([]biolog.Species, error) {
	stmt := `SELECT * FROM species`
	sps := []biolog.Species{}

	if selErr := s.DB.Select(&sps, stmt); selErr != nil {
		return nil, selErr
	}

	return sps, nil
}

// CreateSpecies shrani podatke o doloceni vrsti v naso bazo in vrne dodeljen id
func (s *SpeciesService) CreateSpecies(sp *biolog.Species) (*biolog.Species, error) {
	stmt := `INSERT INTO species (id, species, kingdom, species_family, species_class, phylum, species_order, genus, scientific_name, canonical_name, conservation_status)
		VALUES (:id, :species, :kingdom, :species_family, :species_class, :phylum, :species_order, :genus, :scientific_name, :canonical_name, :conservation_status)`

	_, err := s.DB.NamedExec(stmt, sp)
	if err != nil {
		return nil, err
	}

	// Vrni novo vrsto s pomocjo napisane metode
	return s.Species(*sp.ID)
}

// UpdateSpecies posodobi vrsto s podanim ID glede na nove (non-nil) podatke podane v sp
func (s *SpeciesService) UpdateSpecies(gbifKey int, sp biolog.Species) error {
	q, args := buildInsertUpdateQuery(buildUpdate, "species", sp)
	args = append(args, gbifKey)

	if _, err := s.DB.Exec(q, args...); err != nil {
		return err
	}

	return nil
}

// DeleteSpecies zbrise doloceno vrsto (ce ni navedena v nobenem izmed opazovanj)
func (s *SpeciesService) DeleteSpecies(gbifKey int) error {
	stmt := `DELETE FROM species WHERE id = $1`

	if _, err := s.DB.Exec(stmt, gbifKey); err != nil {
		return err
	}

	return nil
}

// Observation vrne zapis z dolocenim ID
func (s *SpeciesService) Observation(id int) (*biolog.Observation, error) {
	stmt := `SELECT * FROM observation WHERE id = $1`
	ob := &biolog.Observation{}

	if getErr := s.DB.Get(ob, stmt, id); getErr != nil {
		return nil, getErr
	}

	return ob, nil
}

// Observations vrne vse podane zapise o opazenih vrstah (vse, ki so javni)
// TODO:
// 	- preveri za override nad public_observations pri User
// FIXME:
// 	- vracanje lokacije kot koordinate, comma separated (trenutno je HEX)
func (s *SpeciesService) Observations() ([]biolog.Observation, error) {
	stmt := `SELECT * FROM observation WHERE public_visibility = TRUE`
	obs := []biolog.Observation{}

	if selErr := s.DB.Select(&obs, stmt); selErr != nil {
		return nil, selErr
	}

	return obs, nil
}

// CreateObservation kreira nov zapis o opazeni vrsti
func (s *SpeciesService) CreateObservation(o *biolog.Observation) (*biolog.Observation, error) {
	ob := biolog.Observation{}

	q, args := buildInsertUpdateQuery(buildInsert, "observation", o)
	if getErr := s.DB.Get(&ob, q, args...); getErr != nil {
		return nil, getErr
	}

	return &ob, nil
}

// DeleteObservation zbrise dolocen zapis o opazeni vrsti
func (s *SpeciesService) DeleteObservation(id int) error {
	stmt := `DELETE FROM observation WHERE id = $1`

	_, err := s.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateObservation posodobi opazovalni list, ki ima enak ID
// Nove podatke preberemo iz slovarja, pri cemer so kljuci enaki imenom atributov
func (s *SpeciesService) UpdateObservation(id int, ob biolog.Observation) error {
	q, args := buildInsertUpdateQuery(buildUpdate, "observation", ob)
	// Dodaj ID na konec seznama argumentov za query
	args = append(args, id)

	if _, err := s.DB.Exec(q, args...); err != nil {
		return err
	}

	return nil
}

// ConservationStatus vrne podatke o dolocenem statusu ogrozenosti
func (s *SpeciesService) ConservationStatus(id int) (*biolog.ConservationStatus, error) {
	stmt := `SELECT * FROM conservation_status WHERE id = $1`
	cs := &biolog.ConservationStatus{}

	if getErr := s.DB.Get(cs, stmt, id); getErr != nil {
		return nil, getErr
	}

	return cs, nil
}

// ConservationStatuses vrne vsa mozna stanja ogrozenosti za doloceno vrsto
func (s *SpeciesService) ConservationStatuses() ([]biolog.ConservationStatus, error) {
	stmt := `SELECT * FROM conservation_status`
	css := []biolog.ConservationStatus{}

	if selErr := s.DB.Select(&css, stmt); selErr != nil {
		return nil, selErr
	}

	return css, nil
}
