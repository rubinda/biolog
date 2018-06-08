package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rubinda/biolog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleAuth = 1
	jwtSignKey = []byte(viper.GetString("jwt.key"))
)

// JWTToken je swagger model za parameter.
// Pove, da je na zahtevah potreben JWT Token.
//
// swagger:parameters getUsers
type JWTToken struct {
	// JWT Token potreben za avtorizacijo zahteve
	//
	// in: header
	// required: true
	Authorization string
}

// Handler je kolekcija vseh nasih service handler
type Handler struct {
	UserHandler    *UserHandler
	SpeciesHandler *SpeciesHandler
	OAuthConf      *oauth2.Config
	*chi.Mux
}

// GoogleUser je model za odgovor podatkov, ki jih poslje OAuth na Google
type GoogleUser struct {
	// Google ID od uporabnika
	// Tipicno 22 mestno stevilo
	ID string `json:"sub"`

	// Ime uporabnika
	// Tipicno GivenName + FamilyName
	Name string `json:"name"`

	// Ime uporabnika
	GivenName string `json:"given_name"`

	// Priimek uporabnika
	FamilyName string `json:"family_name"`

	// URL do Google strani uporabnika
	Profile string `json:"profile"`

	// URL do slike
	Picture string `json:"picture"`

	// Email naslov uporabnika
	Email string `json:"email"`

	// Vrednost ali je Email preverjen (?)
	EmailVerified bool `json:"email_verified"`
}

// EmailClaims je nadgradnja standardnega Claims pri JWT
type EmailClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// NewRootHandler ustvari starsa vseh ostalih handlerjev, nosi tudi primarni Router
func NewRootHandler(us biolog.UserService, ss biolog.SpeciesService) *Handler {
	h := &Handler{
		Mux: chi.NewRouter(),
	}

	// Ustvari novo konfiguracijo za Google OAuth2, ClientID in ClientSecret
	// se dodata v cmd/biolog/main.go takoj za to funkcijo
	h.OAuthConf = &oauth2.Config{
		ClientID:     viper.GetString("oauth.google.client-id"),
		ClientSecret: viper.GetString("oauth.google.client-secret"),
		RedirectURL:  "https://127.0.0.1:4000/api/v1/authenticate",
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
		r.Group(func(r chi.Router) {
			r.Use(JWTAuthMiddleware)
			r.Mount("/users", h.UserHandler)
		})

		// Podpoti za endpoint '/species'
		h.SpeciesHandler = NewSpeciesHandler()
		h.SpeciesHandler.SpeciesService = ss
		r.Group(func(r chi.Router) {
			r.Use(JWTAuthMiddleware)
			r.Mount("/species", h.SpeciesHandler)
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

// GetUserEmail iz Context-a zahteve pobere uporabnikov email in ga vrne kot string
func getUserEmail(r *http.Request) string {
	return r.Context().Value(contextEmailKey("userEmail")).(string)
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
		Name:     "originalState",
		Value:    state,
		MaxAge:   60,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, sc)
	url := h.OAuthConf.AuthCodeURL(state)
	// Poslji 303 in preusmeri na ustrezen URL
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// AuthHandler je pot, kamor prispe callback iz strani zunanjega avtentikatorja (Google),
// Preveri ali se uporabnik ponovno prijavlja (shrani podatke), ali pa ce je ze registriran z aplikacijo
//
// swagger:route GET /authenticate login authenticate
//
// Poskrbi za callback pri avtentikaciji z zunanjim ponudnikom (Google)
//
// Responses:
//		401:
//      301:
//
// TODO:
// 	- podatke o uporabniku shrani v PB kot biolog.User
// 	- preveri ali uporabnik z id (email) ze obstaja v PB, potem preskoci novo shranjevanje
//  - (?) shrani Token, RefreshToken in ExpiresAt
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Preveri ali se stanje iz *LoginHandler in v odgovoru ujemata
	sc, err := r.Cookie("originalState")
	if err != nil || sc.Value != r.FormValue("state") {
		// Stanje se ne ujema ali pa je prislo do napake, odgovori z 401
		log.Error(err.Error())
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

	// Shrani prejete podatke v GoogleUser
	var gu GoogleUser
	err = json.NewDecoder(userResponse.Body).Decode(&gu)
	if err != nil {
		log.Error(err)
	}
	log.Info(gu.Email)

	// Preveri ali uporabnik ze obstaja (preko unique emaila), ce ne ga shrani
	var u *biolog.User
	u, err = h.UserHandler.UserService.UserByEmail(gu.Email)

	if err != nil {
		if err.Error() == "Not found" {
			// Uporabnik ni bil najden, torej se prijavlja na novo
			// Iz GoogleUser izgradi biolog.User in ga shrani v PB
			u = &biolog.User{
				ExternalID:           &gu.ID,
				DisplayName:          &gu.Name,
				GivenName:            &gu.GivenName,
				FamilyName:           &gu.FamilyName,
				Email:                &gu.Email,
				Picture:              &gu.Picture,
				ExternalAuthProvider: &googleAuth,
			}
			u, err = h.UserHandler.UserService.CreateUser(*u)
			if err != nil {
				log.Error("Uporabnika ni bilo mogoce kreirati")
				log.Error(err)
				respondWithError(w, http.StatusInternalServerError, "Napaka pri kreiranju uporabnika")
			}
		} else {
			log.Error(err.Error())
			respondWithError(w, http.StatusInternalServerError, "Neznana napaka pri kreiranju uporabnika")
		}
	}

	// Dodeli nov JWT uporabniku
	// TODO:
	//  - preveri cas za potek JWT tokena (10-15min ?)
	claims := &EmailClaims{
		gu.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			Issuer:    "biolog-app",
		},
	}
	jwttok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := jwttok.SignedString(jwtSignKey)
	ssJSON, _ := json.Marshal(map[string]string{"token": ss})

	// Odgovori z JWT v telesu zahtevka
	w.WriteHeader(http.StatusOK)
	w.Write(ssJSON)
	// TODO:
	//  - logika za preusmeritev, ali naj bo to na frontend (vrni JWT v Cookie in preusmeri?)
	//  - refresh token
	// glej https://stackoverflow.com/questions/43090518/how-to-properly-handle-a-jwt-refresh

	// Po uspesni prijavi uporabnika preusmeri na domaco stran
	//http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// Za potrebe Context pri JWTAuthMiddleware (se izogne napaki 'Can't use basic type string for context key')
type contextEmailKey string

func (c contextEmailKey) String() string {
	return string(c)
}

// JWTAuthMiddleware se uporabi, da preveri ali ima zahtevek ustrezen JWT in
// mu je dovoljen dostop do vira
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Token loci od polja 'Bearer ' in ga sparsaj
		reqAuth := r.Header.Get("Authorization")
		tokStr := strings.Split(reqAuth, "Bearer ")[1]
		token, err := jwt.ParseWithClaims(tokStr, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.key")), nil
		})

		if claims, ok := token.Claims.(*EmailClaims); ok && token.Valid {
			// Token je veljaven, prav tako smo iz Claims pridobili Email uporabnika ki prozi zahtevo
			if claims.Email == "" {
				respondWithError(w, http.StatusBadRequest, "Tokec nima polja email")
				return
			}
			var emailKey = contextEmailKey("userEmail")
			ctx := context.WithValue(r.Context(), emailKey, claims.Email)

			// Uporabniku dovoli prehod naprej
			next.ServeHTTP(w, r.WithContext(ctx))
		} else if ve, ok := err.(*jwt.ValidationError); ok {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// Token ni pravilne oblike
				respondWithError(w, http.StatusBadRequest, "Tokec ni veljavne oblike")

			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token je bodisi potekel, ali pa se ni veljaven
				respondWithError(w, http.StatusBadRequest, "Tokec vam je potekel")

			} else if ve.Errors&(jwt.ValidationErrorSignatureInvalid) != 0 {
				// Token nima veljavnega podpisa (nekdo ga je spreminjal)
				respondWithError(w, http.StatusBadRequest, "Tokec nima veljavnega podpisa")

			} else {
				log.Info("Something is wrong with the JWT token:", err)
				respondWithError(w, http.StatusBadRequest, "Napaka pri obdelavi tokeca")
			}
		} else {
			log.Info("Couldn't handle this JWT token:", err)
			respondWithError(w, http.StatusBadRequest, "Napaka pri obdelavi tokeca")
		}
	})
}
