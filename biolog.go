// Package biolog only contains simple data types like a User struct for
// holding user data or a UserService interface for fetching or saving user data.
package biolog

import (
	"time"
)

// UserService nudi interface vseh metod za delo z uporabniki
type UserService interface {
	User(id int) (*User, error)
	Users() ([]User, error)
	CreateUser(u User) (*User, error)
	DeleteUser(id int) (int64, error)
	UpdateUser(id int, u User) error
	UserByExtID(id string) (*User, error)

	AuthProvider(id int) (*AuthProvider, error)
	AuthProviders() ([]AuthProvider, error)
}

// User (uporabnik nase aplikacije)
//
// Podatki o uporabniku so vzeti od zunanjega avtentikatorja
//
// swagger:model user
type User struct {
	// 8 mestni ID dolocenega uporabnika
	//
	// required: true
	// min: 10000000
	// max: 99999999
	// example: 10000000
	ID *int `json:"id"`

	// Pove ali uporabnik dovoli javen dostop do svojih opazanj
	PublicObservations *bool `db:"public_observations" json:"publicObservations"`

	// Prikazno ime za uporabnika
	//
	// required: true
	// max length: 64
	// example: David Rubin
	DisplayName *string `db:"display_name" json:"displayName"`

	// ID zunanjega avtentikatorja za tega uporabnika
	//
	// required: true
	// max length: 255
	// example: 8457232358972358923566
	ExternalID *string `db:"external_id" json:"externalID"`

	// Ime uporabnika
	//
	// required: true
	// pattern: [A-Za-z]+
	// max length: 32
	// example: David
	GivenName *string `db:"given_name" json:"givenName"`

	// Priimek uporabnika
	//
	// required: true
	// pattern: [A-Za-z]+
	// max length: 32
	// example: Rubin
	FamilyName *string `db:"family_name" json:"familyName"`

	// Elektronski naslov uporabnika
	//
	// required: true
	// max length: 128
	// swagger:strfmt email
	// example: david@biologapp.com
	Email *string `json:"email"`

	// URL do prikazne slike uporabnika
	//
	// max length: 255
	// example: https://www.biologapp.com/static/img/david-profile.png
	Picture *string `json:"picture,omitempty"`

	// Podatki o zunanjem avtentikatorju uporabnika
	// example: 1
	ExternalAuthProvider *int `db:"external_auth_provider" json:"-"`
}

// AuthProvider (zunanji avtentikator)
//
// Placeholder, ce bi uporabili vec avtentikatorjev (recimo Google + Facebook)
//
// swagger:model authProvider
type AuthProvider struct {
	// Identifikator ponudnika avtentikacije
	//
	// required: true
	// example: 1
	ID int

	// Ime ponudnika avtentikacije
	//
	// required: true
	// max length: 32
	// example: Google
	Name string
}

// SpeciesService nudi interface za delo z vrstami in zapisi o njih
type SpeciesService interface {
	Species(id int) (*Species, error)
	AllSpecies() ([]Species, error)
	CreateSpecies(sp *Species) (*Species, error)
	UpdateSpecies(gbifKey int, sp Species) error
	DeleteSpecies(gbifKey int) error

	Observation(id int) (*Observation, error)
	Observations() ([]Observation, error)
	CreateObservation(o *Observation) (*Observation, error)
	DeleteObservation(id int) error
	UpdateObservation(id int, ob Observation) error

	ConservationStatus(id int) (*ConservationStatus, error)
	ConservationStatuses() ([]ConservationStatus, error)
}

// ConservationStatus (seznam kratic ogrozenosti vrste)
//
//	Podatki so vnaprej doloceni in sicer 10 statusov
//
// swagger:model conservationStatus
type ConservationStatus struct {
	// Identifikator posameznega statusa
	//
	// required: true
	// example: 5
	ID int `json:"id"`

	// Kratica za status
	//
	// required: true
	// min length: 2
	// max length: 2
	// example: VU
	Acronym string `json:"acronym"`

	// Anglesko ime za status
	//
	// required: true
	// max length: 32
	// example: Vulnerable
	NameEN string `db:"name_en" json:"nameEN"`

	// Slovensko ime za status
	//
	// required: true
	// max length: 32
	// example: Ranljive
	NameSI string `db:"name_si" json:"nameSI"`
}

// Species (lokalno shranjena vrsta)
//
// Podatki so vzeti in spletne strani GBIF, tudi ID
//
// swagger:model species
type Species struct {
	// Enolicni identifikator za vrsto,
	// kljuc je enak kot pri GBIF API
	//
	// required: true
	// example: 5231190
	ID *int `json:"id"`

	// Ime vrste
	//
	// required: true
	// max length: 64
	// example: Passer domesticus
	Species *string `json:"species"`

	// Ime kraljestva za vrsto
	//
	// required: true
	// max length: 64
	// example: Animalia
	Kingdom *string `json:"kingdom"`

	// Ime druzine za vrsto
	//
	// required: true
	// max length: 64
	// example: Passeridae
	Family *string `db:"species_family" json:"family"`

	// Ime razreda za vrsto
	//
	// required: true
	// max length: 64
	// example: Aves
	Class *string `db:"species_class" json:"class"`

	// Ime debla za vrsto
	//
	// required: true
	// max length: 64
	// example: Chordata
	Phylum *string `json:"phylum"`

	// Ime reda za vrsto
	//
	// required: true
	// max length: 64
	// example: Passeriformes
	Order *string `db:"species_order" json:"order"`

	// Ime roda za vrsto
	//
	// required: true
	// max length: 64
	// example: Passer
	Genus *string `json:"genus"`

	// Znanstveno ime za vrsto (latinsko ime, vcasih skupaj z avtorjem)
	//
	// required: true
	// max length: 128
	// example: Passer domesticus (Linnaeus, 1758)
	ScientificName *string `db:"scientific_name" json:"scientificName"`

	// Kanonicno ime za vrsto (samo latinsko ime)
	//
	// required: true
	// max length: 128
	// example: Passer domesticus
	CanonicalName *string `db:"canonical_name" json:"canonicalName"`

	// Stanje ogrozenosti vrste
	//
	// required: true
	// min: 1
	// max: 10
	// example: 8
	ConservationStatus *int `db:"conservation_status" json:"conservationStatus"`
}

// Observation (zapis o opazeni vrsti)
//
// Predstavlja opazovalni list, na katerem je zapisana ena opazena vrsta,
// kolicina osebkov ter cas in lokacija
//
// swagger:model observation
type Observation struct {
	// Identifikator opazovalnega lista
	//
	// required: true
	// example: 1
	ID *int `json:"id"`

	// Casovni posnetek trenutka, ko je vrsta bila opazena
	//
	// required: true
	// swagger:strfmt date-time
	// example: 2018-06-04T11:07:37+00:00
	SightingTime *time.Time `db:"sighting_time" json:"sigthingTime"`

	// Lokacija opazanja, ustvari se tocka v skladu z postGIS geography
	//
	// required: true
	// example: -71.060316, 48.432044
	SightingLocation *string `db:"sighting_location" json:"sightingLocation"`

	// Kolicina osebkov opazenih
	//
	// required: true
	// example: 8
	Quantity *int `json:"quantity"`

	// Vidnost opazanja (javno ali zasebno) za posamezen opazovalni list
	// example: true
	PublicVisibility *bool `db:"public_visibility" json:"publicVisibility"`

	// Uporabnik, ki je kreiral opazovalni list
	//
	// required: true
	// min: 10000000
	// max: 99999999
	// example: 10000000
	User *int `db:"biolog_user" json:"user"`

	// Vrsta, ki je bila opazena
	//
	// required: true
	// example: 5231190
	Species *int `json:"species"`
}
