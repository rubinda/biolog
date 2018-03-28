package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rubinda/biolog"
	log "github.com/sirupsen/logrus"
)

// UserHandler predstavlja http handler za nas UserService, prav tako je na njem
// chi Subrouter za ustrezne endpointe
type UserHandler struct {
	UserService biolog.UserService
	*chi.Mux
}

// NewUserHandler generira vse poti, ki se nanasajo na uporabnike
func NewUserHandler() *UserHandler {
	u := &UserHandler{
		//UserService: us,
		Mux: chi.NewRouter(),
	}

	// Prefix do tukaj je ze /api/v1/users, poti pisemo od tega naprej
	u.Get("/", u.GetUsers)
	u.Post("/", u.CreateUser)
	u.Get("/{id:[0-9]+}", u.GetUserByID)
	u.Patch("/{id:[0-9]+}", u.UpdateUser)
	u.Delete("/{id:[0-9]+}", u.DeleteUser)
	u.Get("/{id:[0-9]+}/external", u.GetUserExternalDetails)

	u.Get("/external/{id:[0-9]+}", u.GetExternalUserByID)
	u.Get("/external/{extID}", u.GetExternalUserByExtID)
	u.Patch("/external/{id:[0-9]+}", u.UpdateExternalUser)

	u.Get("/auth_providers", u.GetAuthProviders)
	u.Get("/auth_providers/{id:[0-9]+}", u.GetAuthProvider)

	return u
}

// GetUserByID vrne podrobnosti o uporabniku s podanim ID
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	id, conErr := strconv.Atoi(chi.URLParam(r, "id"))
	usr, err := u.UserService.User(id)
	log.Info("ID gained=", id)
	if conErr != nil || err != nil {
		fmt.Fprintf(w, "Error occured")
	} else {
		fmt.Fprintf(w, usr.DisplayName)
	}
}

// CreateUser ustvari novega uporabnika
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetUsers vrne vse uporabnike
func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// UpdateUser posodobi podatke o dolocenem uporabniku
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// DeleteUser zbrise dolocenega uporabnika
// (!) Zbrisejo se tudi vsi povezani zapisi (ExternalUser, Observations ...).
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetUserExternalDetails vrne podrobnosti o zunanjem uporabniku, katerega referenciramo
// preko uporabniskega racuna User
func (u *UserHandler) GetUserExternalDetails(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetExternalUserByID pridobi podatke o zunanjem uporabniku (uporabnik od zunanjega avtentikatorja),
// preko nasega ID zanj
func (u *UserHandler) GetExternalUserByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetExternalUserByExtID pridobi podatke o zunanjem uporabniku (uporabnik od zunanjega avtentikatorja),
// preko ID, ki ga ima uporabnik pri zunanjem avtentikatorju
func (u *UserHandler) GetExternalUserByExtID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented 2")
	// TODO implement the method
}

// UpdateExternalUser posodobi podatke o zunanjem uporabniku v nasi podatkovni bazi
// preko nasega ID zanj
func (u *UserHandler) UpdateExternalUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetAuthProviders pridobi in izpise vse shranjene zunanje avtentikatorje
func (u *UserHandler) GetAuthProviders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}

// GetAuthProvider pridobi podrobnosti o posameznem ponudniku avtentikacije
func (u *UserHandler) GetAuthProvider(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Not implemented")
	// TODO implement the method
}
