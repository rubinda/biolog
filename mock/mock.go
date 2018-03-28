// Package mock vsebuje mock implementacije od service, spisan je tako, da si pri
// klicu neke funkcije servica najprej se dolocimo kako naj funkcija izgleda
//
// V uporabo pride pri testiranju http handlerjev, kjer nas zanima le njihovo delovanje
package mock

import (
	"github.com/rubinda/biolog"
)

// UserService predstavlja mock za biolog.UserService
type UserService struct {
	UserFn       func(id int) (*biolog.User, error)
	UsersFn      func() ([]*biolog.User, error)
	CreateUserFn func(u *biolog.User) error
	DeleteUserFn func(id int) error

	ExtUserFn       func(id int) (*biolog.ExternalUser, error)
	CreateExtUserFn func(eu *biolog.ExternalUser) error
	DeleteExtUserFn func(id int) error

	AuthProviderFn func(id int) (*biolog.AuthProvider, error)
}

// User mock za vracanje uporabnika preko ID
func (s *UserService) User(id int) (*biolog.User, error) {
	return s.UserFn(id)
}

// Users mock za vracanje vseh uporabnikov
func (s *UserService) Users() ([]*biolog.User, error) {
	return s.UsersFn()
}

// CreateUser mock za ustvarjanje novega uporabnika
func (s *UserService) CreateUser(u *biolog.User) error {
	return s.CreateUserFn(u)
}

// DeleteUser mock za brisanje uporabnika
func (s *UserService) DeleteUser(id int) error {
	return s.DeleteUserFn(id)
}

// ExtUser mock za vracanje zunanjega uporabnika preko ID
func (s *UserService) ExtUser(id int) (*biolog.ExternalUser, error) {
	return s.ExtUser(id)
}

// CreateExtUser mock za ustvarjanje zunanjega uporabnika
func (s *UserService) CreateExtUser(eu *biolog.ExternalUser) error {
	return s.CreateExtUserFn(eu)
}

// DeleteExtUser mock za brisanje zunanjega uporabnika
func (s *UserService) DeleteExtUser(id int) error {
	return s.DeleteExtUserFn(id)
}

// AuthProvider mock za pridobivanje podrobnosti o ponudniku avtentikacije
func (s *UserService) AuthProvider(id int) (*biolog.AuthProvider, error) {
	return s.AuthProviderFn(id)
}

// SpeciesService predstavlja mock za biolog.SpeciesService
type SpeciesService struct {
	SpeciesFn func(id int) (*biolog.Species, error)

	ObservationFn       func(id int) (*biolog.Observation, error)
	ObservationsFn      func() ([]*biolog.Observation, error)
	CreateObservationFn func(o *biolog.Observation) error
	DeleteObservationFn func(id int) error
	UpdateObservationFn func(id int) error

	ConservationStatusFn   func(id int) (*biolog.ConservationStatus, error)
	ConservationStatusesFn func() ([]*biolog.ConservationStatus, error)
}

// Species mock za vracanje vrste preko ID
func (s *SpeciesService) Species(id int) (*biolog.Species, error) {
	return s.SpeciesFn(id)
}

// Observation mock za vracanje opazovalnega lista preko ID
func (s *SpeciesService) Observation(id int) (*biolog.Observation, error) {
	return s.ObservationFn(id)
}

// Observations mock za vracanje vec opazovalnih listov
func (s *SpeciesService) Observations() ([]*biolog.Observation, error) {
	return s.ObservationsFn()
}

// CreateObservation mock za kreiranje opazovalnega lista
func (s *SpeciesService) CreateObservation(o *biolog.Observation) error {
	return s.CreateObservationFn(o)
}

// DeleteObservation mock za brisanje opazovalnega lista
func (s *SpeciesService) DeleteObservation(id int) error {
	return s.DeleteObservationFn(id)
}

// ConservationStatus vrne podatke o dolocenem statusu ogrozenosti
func (s *SpeciesService) ConservationStatus(id int) (*biolog.ConservationStatus, error) {
	return s.ConservationStatusFn(id)
}

// ConservationStatuses vrne vse mozne statuse ogrozenosti za doloceno vrsto
func (s *SpeciesService) ConservationStatuses() ([]*biolog.ConservationStatus, error) {
	return s.ConservationStatusesFn()
}

// UpdateObservation mock za posodabljanje opazovalnega lista
func (s *SpeciesService) UpdateObservation(id int) error {
	return s.UpdateObservationFn(id)
}
