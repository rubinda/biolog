// Package biolog only contains simple data types like a User struct for
// holding user data or a UserService interface for fetching or saving user data.
package biolog

import "time"

// UserService nudi interface vseh metod za delo z uporabniki
type UserService interface {
	User(id int) (*User, error)
	Users() ([]User, error)
	CreateUser(u User) (*User, error)
	DeleteUser(id int) (int64, error)
	UpdateUser(id int, u User) error
	UserByExtID(id string) (*ExternalUser, error)

	AuthProvider(id int) (*AuthProvider, error)
	AuthProviders() ([]AuthProvider, error)
}

// User je posplosen model uporabnika
// Pri 'PublicObservations' je v golang default enak false
type User struct {
	ID                 *int
	PublicObservations *bool   `db:"public_observations"`
	DisplayName        *string `db:"display_name"`
}

// ExternalUser deduje od User in predstavlja podatke pridobljene
// iz strani zunanjega avtentikatorja
type ExternalUser struct {
	ID                   *int
	ExternalID           *string `db:"external_id"`
	GivenName            *string `db:"given_name"`
	FamilyName           *string `db:"family_name"`
	Email                *string
	Picture              *string
	ExternalAuthProvider *int `db:"external_auth_provider"`
	User                 *int `db:"biolog_user"`
}

// AuthProvider je zunanji avtentikator za prijavo v aplikacijo
type AuthProvider struct {
	ID   int
	Name string
}

// SpeciesService nudi interface za delo z vrstami in zapisi o njih
type SpeciesService interface {
	Species(id int) (*Species, error)
	CreateSpecies(sp *Species) (int64, error)

	Observation(id int) (*Observation, error)
	Observations() ([]*Observation, error)
	CreateObservation(o *Observation) (int, error)
	DeleteObservation(id int) error
	UpdateObservation(o map[string]string) error

	ConservationStatus(id int) (*ConservationStatus, error)
	ConservationStatuses() ([]*ConservationStatus, error)
}

// ConservationStatus je seznam kratic ogrozenosti vrste
type ConservationStatus struct {
	ID      int
	Acronym string
	NameEN  string `db:"name_en"`
	NameSI  string `db:"name_si"`
}

// Species je vrsta, ki je bila opazena in zabelezena
type Species struct {
	ID                 int
	Species            string
	Kingdom            string
	Family             string
	Class              string
	Phylum             string
	Order              string
	Genus              string
	ScientificName     string `db:"scientific_name"`
	CanonicalName      string `db:"canonical_name"`
	ConservationStatus int    `db:"conservation_status"`
	GBIFKey            int    `db:"gbif_key"`
}

// Observation je zapis, da je na podani lokaciji bila opazena vrsta
// Pri 'PublicVisibility' je v golang default enak false
type Observation struct {
	ID               int
	SightingTime     time.Time `db:"sighting_time"`
	SightingLocation string    `db:"sighting_location"`
	Quantity         int
	PublicVisibility bool `db:"public_visibility"`
	User             int
	Species          int
}
