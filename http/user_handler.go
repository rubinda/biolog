package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rubinda/biolog"

	log "github.com/sirupsen/logrus"
	//	log "github.com/sirupsen/logrus"
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
// TODO:
// 	- vrnejo se naj le uporabniki, ki imajo vsaj 1 javno opazanje
// 	- boljse javljanje napak
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r)
	if parseErr {
		return
	}

	// Pridobi uporabnika preko baze
	usr, err := u.UserService.User(id)

	// Preveri napake pri pridobivanju iz PB in ustrezno obvesti odjemalca
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, usr)
}

// CreateUser ustvari novega uporabnika
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var usr biolog.User

	// Pridobi podatke o novem uporabniku iz telesa zahtevka
	decErr := json.NewDecoder(r.Body).Decode(&usr)
	if decErr != nil {
		switch {
		case decErr == io.EOF:
			respondWithError(w, 400, "Telo zahtevka pri kreiranju uporabnika ne more biti prazno")
		default:
			respondWithError(w, 400, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}
	// Kreiraj novega uporabnika
	newUsr, err := u.UserService.CreateUser(&usr)

	// Napaka pri kreiranju
	if err != nil {
		respondWithError(w, 400, "Napaka pri ustvarjanju novega uporabnika")
		log.Error(err)
		return
	}

	// Uporabnik uspesno ustvarjen, poslji nazaj na novo shranjene podatke
	respondWithJSON(w, 201, newUsr)
}

// GetUsers vrne vse uporabnike
// TODO:
// 	- vrnejo se naj le uporabniki, ki imajo javna opazanja
// 	- paginacija
func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var usrs []biolog.User

	// Pridobi podatke o vseh uporabnikih, ki so na voljo
	usrs, err := u.UserService.Users()

	// Preveri ali je prislo do napake
	if err != nil {
		respondWithError(w, 400, "Pri poizvedbi nad vsemi uporabniki je prislo do napake")
		return
	}

	// Odgovori s seznamom vseh uporabnikov
	respondWithJSON(w, 200, usrs)
}

// UpdateUser posodobi podatke o dolocenem uporabniku
// FIXME:
// 	- branje ID iz telesa in ID iz URL
// TODO:
// 	- javljanje napak (neveljavni znaki za polja?)
//  - uporabnik lahko posodablja le lasten racun
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r)
	if parseErr {
		return
	}

	var usr biolog.User
	// Pridobi podatke o uporabniku iz telesa zahtevka
	decErr := json.NewDecoder(r.Body).Decode(&usr)
	if decErr != nil {
		switch {
		case decErr == io.EOF:
			respondWithError(w, 400, "Telo zahtevka pri kreiranju uporabnika ne more biti prazno")
		default:
			respondWithError(w, 400, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}
	usr.ID = &id
	if updErr := u.UserService.UpdateUser(id, usr); updErr != nil {
		respondWithError(w, 400, updErr.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}

// DeleteUser zbrise dolocenega uporabnika
// (!) Zbrisejo se tudi vsi povezani zapisi (ExternalUser, Observations ...).
// TODO:
// 	- preveri da uporabnik lahko zbrise le svoj racun
//  - javljanje napak (unauthorized, non-existent)
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r)
	if parseErr {
		return
	}

	_, err := u.UserService.DeleteUser(id)

	// Preveri ce je prislo do napake
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// Uporabnik uspesno zbrisan, poslji 204
	respondWithJSON(w, 204, nil)

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
	// TODO: implement the method
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

// GetIDFromURL pridobi veljaven url iz poti v URL, ce pride do napake obvesti odjemalca
// Vraca veljaven id (stevilo int32) in vrednost ali je prislo do napake
func getIDFromURL(w http.ResponseWriter, r *http.Request) (int, bool) {
	// Pridobi ID iz URL in ga pretvori v stevilo (Router poskrbi da je na tej poti vedno stevilka,
	// zato lahko napako ignoriramo)
	id64, parErr := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	id := int(id64)

	// Pri pretvarjanju v Integer (od 0 do 2^31 -1) je prislo do napake (najverjetneje je overflow)
	if parErr != nil {
		e := parErr.(*strconv.NumError)
		// Obvesti, da ID ni v veljavnem obsegu
		if e.Err == strconv.ErrRange {
			respondWithError(w, 400, "Neveljaven ID: izven obsega")
			// Prislo je do druge napake
		} else {
			respondWithError(w, 400, parErr.Error())
		}
		return 0, true
	}

	return id, false
}
