package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rubinda/biolog"
	//	log "github.com/sirupsen/logrus"
)

// UserHandler predstavlja http handler za nas UserService, prav tako je na njem
// chi Subrouter za ustrezne endpointe
type UserHandler struct {
	UserService biolog.UserService
	*chi.Mux
}

// UserID parameter model.
//
// Uporablja se za operacije, ki pricakujejo ID uporabnika v poti
// swagger:parameters getUserByID deleteUser updateUser
type UserID struct {
	// ID uporabnika
	//
	// in: path
	// required: true
	// min: 10000000
	// max: 99999999
	ID int32
}

// UserExtID parameter model.
//
// Uporablja se za operacije, ki se sklicujejo na uporabnike glede na id zunajnega avtentikatorja
// swagger:parameters getUserByExtID
type UserExtID struct {
	// ID zunanjega avtentikatorja za uporabnika
	//
	// in: path
	// required: true
	extID int
}

// UserBodyParams model.
//
// Uporablja se pri operacijah, ki zahtevajo podatke o uporabniku v
// telesu zahtevka
// swagger:parameters createUser
type UserBodyParams struct {
	// Podatki o uporabniku'
	//
	// in: body
	// required: true
	User *biolog.User
}

// AuthProviderResponse model.
//
// Se uporablja pri virih, ki vracajo podatke o zunanjih avtentikatorjih
// swagger:response authProviderResponse
type AuthProviderResponse struct {
	// in: body
	Payload *biolog.AuthProvider `json:"authProvider"`
}

// AuthProvidersID parameter model
//
// Uporablja se za operacijo, ko pridobivamo podrobnosti o
// ponudniku avtentikacije
// swagger:parameters getAuthProviderByID
type AuthProvidersID struct {
	// ID ponudnika avtentikacije
	//
	// in: path
	// required: true
	// min: 1
	ID int
}

// NewUserHandler generira vse poti, ki se nanasajo na uporabnike
func NewUserHandler() *UserHandler {
	u := &UserHandler{
		//UserService: us,
		Mux: chi.NewRouter(),
	}
	// Prefix do tukaj je ze /api/v1/users, poti pisemo od tega naprej

	// swagger:route GET /users users getUsers
	//
	// Pridobi vse uporabnike, ki imajo vsaj 1 javno opazanje
	//
	// Responses:
	//		400: description: Prislo je do napake
	//		200: []user
	u.Get("/", u.GetUsers)

	// TODO:
	//	- pridobi ID iz URL preko middleware
	u.Route("/{id:\\d{8}}", func(r chi.Router) {
		// swagger:route GET /users/{id} users getUserByID
		//
		// Pridobi podrobnosti o uporabniku
		//
		// Responses:
		//		400: description: Prislo je do napake
		//		200: user
		r.Get("/", u.GetUserByID)

		// swagger:route PATCH /users/{id} users updateUser
		//
		// Posodobi podatke o uporabniku
		//
		// Responses:
		//		400: description: Prislo je do napake
		// 		204:
		r.Patch("/", u.UpdateUser)

		// swagger:route DELETE /users/{id} users deleteUser
		//
		// Zbrise uproabniski racun in vse zapise o uporabniku
		//
		// Responses:
		//		400: description: Prislo je do napake
		//		204:
		r.Delete("/", u.DeleteUser)
	})

	// Metode za ponudnike zunanje avtentikacije

	// swagger:route GET /auth_providers authproviders getAuthProviders
	//
	// Pridobi vse mozne ponudnike avtentikacije
	//
	// Responses:
	//		400: description: Prislo je do napake
	//		200: []authProviderResponse
	u.Get("/auth_providers", u.GetAuthProviders)

	// swagger:route GET /auth_providers/{id} authproviders getAuthProviderByID
	//
	// Pridobi podrobnosti o dolocenem ponudniku avtentikacije, ki je pri nas na voljo
	//
	// Responses:
	//		400: description: Prislo je do napake
	// 		200: authProviderResponse
	u.Get("/auth_providers/{id:[0-9]+}", u.GetAuthProvider)

	return u
}

// GetUserByID vrne podrobnosti o uporabniku s podanim ID
// TODO:
// 	- boljse javljanje napak
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	// Pridobi uporabnika preko baze
	usr, err := u.UserService.User(id)

	// Preveri napake pri pridobivanju iz PB in ustrezno obvesti odjemalca
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, usr)
}

// GetUsers vrne vse uporabnike
// TODO:
// 	- vrnejo se naj le uporabniki, ki imajo javna opazanja
// 	- paginacija
func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Pridobi podatke o vseh uporabnikih, ki so na voljo
	usrs, err := u.UserService.Users()

	// Preveri ali je prislo do napake
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Pri poizvedbi nad vsemi uporabniki je prislo do napake")
		return
	}

	// Odgovori s seznamom vseh uporabnikov
	respondWithJSON(w, http.StatusOK, usrs)
}

// UpdateUser posodobi podatke o dolocenem uporabniku
// FIXME:
// 	- branje ID iz telesa in ID iz URL
// TODO:
// 	- javljanje napak (neveljavni znaki za polja?)
//  - uporabnik lahko posodablja le lasten racun
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	var usr biolog.User
	// Pridobi podatke o uporabniku iz telesa zahtevka
	decErr := json.NewDecoder(r.Body).Decode(&usr)
	if decErr != nil {
		switch {
		case decErr == io.EOF:
			respondWithError(w, http.StatusBadRequest, "Telo zahtevka pri kreiranju uporabnika ne more biti prazno")
		default:
			respondWithError(w, http.StatusBadRequest, "Napaka pri pretvarjanju JSONa iz telesa zahtevka")
		}
		return
	}
	usr.ID = &id
	if updErr := u.UserService.UpdateUser(id, usr); updErr != nil {
		respondWithError(w, http.StatusBadRequest, updErr.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// DeleteUser zbrise dolocenega uporabnika
// (!) Zbrisejo se tudi vsi povezani zapisi (ExternalUser, Observations ...).
// TODO:
//	- zbrise se naj se ExternalUser
// 	- preveri da uporabnik lahko zbrise le svoj racun
//  - javljanje napak (unauthorized, non-existent)
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	_, err := u.UserService.DeleteUser(id)

	// Preveri ce je prislo do napake
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Uporabnik uspesno zbrisan, poslji http.StatusNoContent
	respondWithJSON(w, http.StatusNoContent, nil)

}

// GetAuthProviders pridobi in izpise vse shranjene zunanje avtentikatorje
func (u *UserHandler) GetAuthProviders(w http.ResponseWriter, r *http.Request) {
	// Pridobi podatke o vseh ponudnikih avtentikacije
	ps, err := u.UserService.AuthProviders()

	// Preveri ali je prislo do napake
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Pri pridobivanju ponudnikov avtentikacije je prislo do napake")
		return
	}

	// Odgovori s seznamom vseh uporabnikov
	respondWithJSON(w, http.StatusOK, ps)
}

// GetAuthProvider pridobi podrobnosti o posameznem ponudniku avtentikacije
func (u *UserHandler) GetAuthProvider(w http.ResponseWriter, r *http.Request) {
	id, parseErr := getIDFromURL(w, r, "id")
	if parseErr {
		return
	}

	p, err := u.UserService.AuthProvider(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}
