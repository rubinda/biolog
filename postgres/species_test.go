package postgres_test

import (
	"testing"
	"time"

	"github.com/rubinda/biolog"
	"github.com/stretchr/testify/assert"
)

// TestSpecies preveri vracanje vrste glede na lokalen ID
func TestSpecies(t *testing.T) {
	id := 1
	species, getErr := speciesServiceTest.Species(id)
	if assert.NoError(t, getErr) {
		actualSpecies := biolog.Species{}
		selectErr := speciesServiceTest.DB.Get(&actualSpecies,
			"SELECT * FROM species WHERE ID = $1 LIMIT 1", id)
		if assert.NoError(t, selectErr) {
			assert.Equal(t, actualSpecies, species)
		}
	}
}

// TestCreateSpecies preveri shranjevanje podatkov o neki vrsti v naso bazo
// Za preverjanje se uporabijo podatki pridobljeni s spletne strani GBIf Species API
// (?) Mogoce assert.Equal zajamra, ker 'cases' nimajo IDjev, morebitna resitev
// (?) 	bi lahko bila c.Species.ID = newID
func TestCreateSpecies(t *testing.T) {
	cases := []struct {
		Species *biolog.Species
	}{
		{
			Species: &biolog.Species{Species: "Ailurus fulgens", Kingdom: "Animalia", Family: "Ailuridae",
				Class: "Mammalia", Phylum: "Chordata", Order: "Carnivora", Genus: "Ailurus",
				ScientificName: "Ailurus fulgens Geoffroy Saint-Hilaire & Cuvier, 1825",
				CanonicalName:  "Ailurus fulgens", ConservationStatus: 4, GBIFKey: 5219446},
		}, {
			Species: &biolog.Species{Species: "Hirundo rustica", Kingdom: "Animalia", Family: "Hirundinidae",
				Class: "Aves", Phylum: "Chordata", Order: "Passeriformes", Genus: "Hirundo",
				ScientificName: "Hirundo rustica Linnaeus, 1758", CanonicalName: "Hirundo rustica",
				ConservationStatus: 8, GBIFKey: 9515886},
		},
	}

	for _, c := range cases {
		newID, createErr := speciesServiceTest.CreateSpecies(c.Species)
		if assert.NoError(t, createErr) {
			newSpecies := biolog.Species{}
			getErr := speciesServiceTest.DB.Get(&newSpecies,
				"SELECT * FROM species WHERE id = $1 LIMIT 1", newID)
			if assert.NoError(t, getErr) {
				c.Species.ID = int(newID)
				assert.Equal(t, c.Species, newSpecies)
			}
		}
	}
}

// TestObservarion preveri pridobivanje vrste glede na podan ID
func TestObservation(t *testing.T) {
	ID := 1
	o, selectErr := speciesServiceTest.Observation(ID)
	if assert.NoError(t, selectErr) {
		actualO := &biolog.Observation{}
		getErr := speciesServiceTest.DB.Get(actualO, `SELECT * FROM observation WHERE id = $1`, ID)
		if assert.NoError(t, getErr) {
			assert.Equal(t, actualO, o)
		}
	}
}

// TestObservations vrne vsa javna opazanja (Javna opazanja so tista, pri katerih ima uporabnik PublicObservations
// nastavljen na true, prav tako pa posamezno opazanje rabi PublicVisibility enak true)
func TestObservations(t *testing.T) {
	o, getErr := speciesServiceTest.Observations()
	if assert.NoError(t, getErr) {
		actualO := &[]biolog.Observation{}
		selectErr := speciesServiceTest.DB.Select(actualO, `SELECT o.id, o.quantity, ST_AsText(o.sighting_location) as sighting_location, o.sighting_time, o.quantity, o.biolog_user, o.species FROM observation AS o, biolog_user AS bu
		WHERE o.biolog_user = bu.id AND bu.public_observations = TRUE AND o.public_visibility = TRUE`)
		if assert.NoError(t, selectErr) {
			assert.Equal(t, actualO, o)
		}
	}
}

// TestCreateObservation preveri kreiranje zapisa o opazanju vrste
func TestCreateObservation(t *testing.T) {
	cases := []struct {
		Observation *biolog.Observation
	}{
		{
			Observation: &biolog.Observation{SightingTime: time.Now(), SightingLocation: "46.33061, 15.48705", Quantity: 3,
				PublicVisibility: true, User: 1, Species: 1},
		},
	}
	for _, c := range cases {
		newID, createErr := speciesServiceTest.CreateObservation(c.Observation)
		if assert.NoError(t, createErr) {
			c.Observation.ID = newID
			actualObservation := &biolog.Observation{}
			getErr := speciesServiceTest.DB.Get(actualObservation,
				`SELECT id, sighting_time, ST_AsLatLonText(ST_AsText(sighting_location)) AS sighting_location, quantity, 
				biolog_user, species FROM observation WHERE ID = $1`, newID)
			if assert.NoError(t, getErr) {
				assert.Equal(t, c.Observation, actualObservation)
			}
		}
	}
}

// TestDeleteObservation preveri brisanje dolocenega zapisa o opazanju
func TestDeleteObservation(t *testing.T) {
	ID := 1
	delErr := speciesServiceTest.DeleteObservation(ID)
	if assert.NoError(t, delErr) {
		var oExists bool
		getErr := speciesServiceTest.DB.Get(oExists, `SELECT EXISTS (SELECT 1 FROM observation WHERE id = $1)`, ID)
		if assert.NoError(t, getErr) {
			assert.Equal(t, false, oExists)
		}
	}
}
