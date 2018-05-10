package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Dodatek za PostgreSQL
	"github.com/rubinda/biolog"
)

// UserService predstavlja PostgreSQL implementacijo od biolog.UserService
type UserService struct {
	DB *sqlx.DB
}

// User vrne uporabnika, ki pripada podanemu ID
// FIXME:
// 	- pripadajoci test
func (s *UserService) User(id int) (*biolog.User, error) {
	stmt := `SELECT * FROM biolog_user WHERE id = $1`
	u := &biolog.User{}
	if getErr := s.DB.Get(u, stmt, id); getErr != nil {
		if getErr == sql.ErrNoRows {
			return nil, errors.New("Uporabnik s tem ID ne obstaja")
		}
		return nil, getErr
	}
	return u, nil
}

// Users vrne vse uporabnike
// FIXME:
// 	- pripadajoci test
func (s *UserService) Users() ([]biolog.User, error) {
	stmt := `SELECT * FROM biolog_user`
	us := []biolog.User{}
	if getErr := s.DB.Select(&us, stmt); getErr != nil {
		return nil, getErr
	}
	return us, nil
}

// CreateUser ustvari novega uporabnika za uporabo aplikacije
// FIXME:
// 	- pripadajoci test
func (s *UserService) CreateUser(u biolog.User) (*biolog.User, error) {
	newUser := biolog.User{}

	// Po koncani kreaciji naj se vrne nov dodeljen zapis o uporabniku
	q, args := buildInsertUpdateQuery(buildInsert, "biolog_user", u)
	if err := s.DB.Get(&newUser, q, args...); err != nil {
		return nil, err
	}

	return &newUser, nil
}

// DeleteUser izbrise podanega uporabnika iz podatkovne baze. Javi napako, ce ima uporabnik zapise o opazanjih.
// TODO:
// 	- dodaj Cascade, ki zbrise se vse povezane zapise
// FIXME:
// 	- pripadajoci test
func (s *UserService) DeleteUser(id int) (int64, error) {
	deleteUser := `DELETE FROM biolog_user WHERE ID = $1`
	result, createErr := s.DB.Exec(deleteUser, id)
	if createErr != nil {
		return -1, createErr
	}
	rowsDeleted, _ := result.RowsAffected()
	return rowsDeleted, nil
}

// UpdateUser delno posodobi podatke o uporabniku
//
// (?) Ali je lahko sporno da posodabljas podatke, ki so pridobljeni od zunanjega avtentikatorja?
// FIXME:
// 	- pripadajoci test
func (s *UserService) UpdateUser(id int, u biolog.User) error {
	query, args := buildInsertUpdateQuery(buildUpdate, "biolog_user", u)
	// Dodaj ID v seznam argumentov
	args = append(args, id)
	if _, err := s.DB.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

// UserByExtID vrne zunanjega uporabnika glede na ID zunanjega avtentikatorja
func (s *UserService) UserByExtID(id string) (*biolog.User, error) {
	stmt := `SELECT * FROM biolog_user WHERE external_id = $1`
	eu := &biolog.User{}

	// Pozene poizvedbo in preveri za napake
	if err := s.DB.Get(eu, stmt, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Uporabnika s tem ID ni mogoče najti")
		}
		return nil, err
	}
	return eu, nil
}

// AuthProvider vrne podrobnosti o dolocenem ponudniku avtentikacije
func (s *UserService) AuthProvider(id int) (*biolog.AuthProvider, error) {
	stmt := `SELECT * FROM external_auth_provider WHERE id = $1`
	var authPro biolog.AuthProvider

	if err := s.DB.Get(&authPro, stmt, id); err != nil {
		return nil, err
	}

	return &authPro, nil
}

// AuthProviders vrne vse podatke o vseh zunanjih avtentikatorjih
func (s *UserService) AuthProviders() ([]biolog.AuthProvider, error) {
	stmt := `SELECT * FROM external_auth_provider`
	var authPros []biolog.AuthProvider

	if err := s.DB.Select(&authPros, stmt); err != nil {
		return nil, err
	}

	return authPros, nil
}
