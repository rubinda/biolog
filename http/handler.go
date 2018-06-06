package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rubinda/biolog"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Handler je kolekcija vseh nasih service handler
type Handler struct {
	UserHandler    *UserHandler
	SpeciesHandler *SpeciesHandler
	OAuthConf      *oauth2.Config
	*chi.Mux
}

// NewRootHandler ustvari starsa vseh ostalih handlerjev, nosi tudi primarni Router
func NewRootHandler(us biolog.UserService, ss biolog.SpeciesService) *Handler {
	h := &Handler{
		Mux: chi.NewRouter(),
	}

	// Ustvari novo konfiguracijo za Google OAuth2, ClientID in ClientSecret
	// se dodata v cmd/biolog/main.go takoj za to funkcijo
	h.OAuthConf = &oauth2.Config{
		RedirectURL: "https://127.0.0.1:4000/api/v1/authenticate",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	// A good base middleware stack
	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	// Timeout na zahteve
	h.Use(middleware.Timeout(60 * time.Second))
	// Nastavimo predpono za api
	h.Route("/api/v1", func(r chi.Router) {

		// Podpoti za endpoint '/users'
		h.UserHandler = NewUserHandler()
		h.UserHandler.UserService = us
		// Ustvari nov router z 'fresh middleware stack'
		r.Group(func(ru chi.Router) {
			ru.Mount("/users", h.UserHandler)
		})

		// Podpoti za endpoint '/species'
		h.SpeciesHandler = NewSpeciesHandler()
		h.SpeciesHandler.SpeciesService = ss
		r.Group(func(rs chi.Router) {
			rs.Mount("/species", h.SpeciesHandler)
		})

		// Podpoti za preusmeranje prijav na ponudnika avtentikacije
		r.Route("/login", func(r chi.Router) {

			r.Get("/google", h.GoogleLoginHandler)
		})

		// Podpoti za callback od ponudnikov avtentikacije
		r.Route("/authenticate", func(r chi.Router) {
			r.Get("/", h.AuthHandler)
		})

	})

	return h
}

// RespondWithError vrne napako kot odgovor na http request s podanimi podrobnostmi
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON vrne JSON kot odgovor na zahtevo. Parametra sta http koda odgovora in telo
// FIXME:
// 	- moznost dodajanja lastnih headerjev
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetIDFromURL pridobi veljavno 32 bitno stevilo, pri cemer podamo ime parametra, ce pride do napake obvesti odjemalca
// Vraca veljaven id (stevilo int32) in vrednost ali je prislo do napake
// TODO:
//	- spremeni funkcijo v middleware z Context in uporabi router.Use(GetUser), kar nam bo pridobilo uporabnika v Context
func getIDFromURL(w http.ResponseWriter, r *http.Request, parameter string) (int, bool) {
	// Pridobi ID iz URL in ga pretvori v stevilo (Router poskrbi da je na tej poti vedno stevilka,
	// zato lahko napako ignoriramo)
	id64, parErr := strconv.ParseInt(chi.URLParam(r, parameter), 10, 32)
	id := int(id64)

	// Pri pretvarjanju v Integer (od 0 do 2^31 -1) je prislo do napake (najverjetneje je overflow)
	if parErr != nil {
		e := parErr.(*strconv.NumError)
		// Obvesti, da ID ni v veljavnem obsegu
		if e.Err == strconv.ErrRange {
			respondWithError(w, http.StatusBadRequest, "Neveljaven ID: izven obsega")
			// Prislo je do druge napake
		} else {
			respondWithError(w, http.StatusBadRequest, parErr.Error())
		}
		return 0, true
	}

	return id, false
}

// Vzeto iz https://skarlso.github.io/2016/06/12/google-signin-with-go/,
// preveri za state odgovora in zahteve, kar zasciti pred CSRF napadi
func (h *Handler) getLoginURL(state string) string {
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return h.OAuthConf.AuthCodeURL(state)
}

// RandToken kreira nakljucen 32 znakov dolg niz, ki ga uporabimo za stanje (state) pri Google Auth
// Stanje pomaga pri preprecevanju CSRF napadov
//
// Vzeto iz https://skarlso.github.io/2016/06/12/google-signin-with-go/
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Se izogne opozorilo, naj se osnovni tip string ne uporabi kot Context key
type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// GoogleLoginHandler preusmeri prihajajoco zahtevo na ustrezen Google API
// URL na katerem se uporabnik vpise
//
// swagger:route GET /login/google login loginGoogle
//
// Preusmeri povezavo na ustrezen Google login
//
// Responses:
// 		303:
func (h *Handler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Ustavi nakljucno zaporedje znakov za stanje in ga dodaj v Cookie na povezavo
	// Cookie je veljaven samo 60 sekund
	state := randToken()
	sc := &http.Cookie{
		Name:   "originalState",
		Value:  state,
		MaxAge: 60,
		Path:   "/",
	}
	http.SetCookie(w, sc)
	url := h.OAuthConf.AuthCodeURL(state)
	// Poslji 303 in preusmeri na ustrezen URL
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// AuthHandler je pot, kamor prispe preusmerjena povezava iz strani zunanjega avtentikatorja (Google),
// Preveri ali se uporabnik ponovno prijavlja (shrani podatke), ali pa ce je ze registriran z aplikacijo
//
// swagger:route GET /authenticate login authenticate
//
// Poskrbi za callback pri avtentikaciji z zunanjim ponudnikom (Google)
//
// Responses:
//		401:
//
// TODO:
// 	- podatke o uporabniku shrani v PB kot biolog.User
// 	- preveri ali uporabnik z id (email) ze obstaja v PB
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Preveri ali se stanje iz *LoginHandler in v odgovoru ujemata
	sc, err := r.Cookie("originalState")
	if err != nil || sc.Value != r.FormValue("state") {
		// Stanje se ne ujema ali pa je prislo do napake, odgovori z 401
		http.Error(w, "Neveljavno stanje v odgovoru", http.StatusUnauthorized)
		return
	}

	// Zamenjaj avtorizacijsko kodo pridobljeno iz prvotne preusmeritve za Token, s katerim lahko pridobimo podrobnosti o uporabniku
	tok, err := h.OAuthConf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Preveri ali je token veljaven
	if tok.Valid() == false {
		http.Error(w, "Tokec je neveljaven", http.StatusUnauthorized)
	}

	// Preko klienta poslji zahtevek s tokenom na naslov za pridobivanje osnovnih podatkov o uporabniku
	client := h.OAuthConf.Client(oauth2.NoContext, tok)
	userResponse, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userData, err := ioutil.ReadAll(userResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(userData))

	// Po uspesni prijavi uporabnika preusmeri na domaco stran
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
