package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Dodatek za PostgreSQL
	"github.com/rubinda/biolog"
	log "github.com/sirupsen/logrus"
)

// UserService predstavlja PostgreSQL implementacijo od biolog.UserService
type UserService struct {
	DB *sqlx.DB
}

// User vrne uporabnika, ki pripada podanemu ID
func (s *UserService) User(id int) (*biolog.User, error) {
	stmt := `SELECT * FROM biolog_user WHERE id = $1`
	u := &biolog.User{}
	if getErr := s.DB.Get(u, stmt, id); getErr != nil {
		if getErr == sql.ErrNoRows {
			return nil, errors.New("Uporabnik ne obstaja")
		}
		return nil, getErr
	}
	return u, nil
}

// Users vrne vse uporabnike
func (s *UserService) Users() ([]biolog.User, error) {
	stmt := `SELECT * FROM biolog_user`
	us := []biolog.User{}
	if getErr := s.DB.Select(&us, stmt); getErr != nil {
		return nil, getErr
	}
	return us, nil
}

// CreateUser ustvari novega uporabnika za uporabo aplikacije
func (s *UserService) CreateUser(u *biolog.User) (*biolog.User, error) {
	newUser := biolog.User{}

	// Po koncani kreaciji naj se vrne nov dodeljen zapis o uporabniku
	insertQuery := `INSERT INTO biolog_user(display_name, public_observations)
		VALUES(:display_name, :public_observations)
		RETURNING *`
	insStmt, stmtErr := s.DB.PrepareNamed(insertQuery)
	defer insStmt.Close()
	if stmtErr != nil {
		return nil, stmtErr
	}

	// Pozeni INSERT stavek in pridobi nazaj novega uporabnika
	runErr := insStmt.Get(&newUser, u)
	if runErr != nil {
		return nil, runErr
	}
	return &newUser, nil
}

// DeleteUser izbrise podanega uporabnika iz podatkovne baze. Javi napako, ce ima uporabnik zapise o opazanjih.
// (?) Dodaj DeleteUserCascade, ki zbrise se vse povezane zapise
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
func (s *UserService) UpdateUser(id int, u biolog.User) error {
	query, args := buildInsertUpdateQuery("UPDATE", "biolog_user", u)
	log.Info("Query = ", query)
	if _, err := s.DB.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

// ExtUser vrne zunanjenga uporabnika s podanim ID
func (s *UserService) ExtUser(id int) (*biolog.ExternalUser, error) {
	return nil, errors.New("Not implemented")
}

// CreateExtUser ustvari nov zapis o zunanjem uporabniku
func (s *UserService) CreateExtUser(eu *biolog.ExternalUser) error {
	return errors.New("Not implemented")
}

// DeleteExtUser zbrise zunanjega uporabnika. Javi napako, ce ima uporabnik zapise o opazanjih.
// (?) Dodaj DeleteExtUserCascade, ki zbrise se vse povezane zapise
func (s *UserService) DeleteExtUser(id int) error {
	return errors.New("Not implemented")
}

// AuthProvider vrne podrobnosti o dolocenem ponudniku avtentikacije
func (s *UserService) AuthProvider(id int) (*biolog.AuthProvider, error) {
	return nil, errors.New("Not implemented")
}

// GetUserHandler je funkcija, ki jo klice router
// Klice funkcijo za pridobivanje podrobnosti o podanem uporabniku (preko ID)
/*func (a *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Poberi spremenljivke iz naslova

	tid, err := strconv.Atoi(vars["id"]) // Pridobi ID iz url
	if err != nil {
		// Prislo je do napake pri pridobivanju ID-ja, odgovori z BAD_REQUEST
		respondWithError(w, http.StatusBadRequest, "Neveljaven ID uporabnika")
		return
	}

	user := User{ID: tid}

	if err := user.Get(a.DB); err != nil {
		// Obvesti uporabnika ustrezno glede na napako
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Uporabnik ni bil najden")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	respondWithJSON(w, http.StatusOK, user)
}*/
