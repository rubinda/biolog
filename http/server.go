package http

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// NewServer kreira nov server na podanem naslovu s podanim handlerjem
func NewServer(addr string, h *Handler) *http.Server {
	// (?) Good practice to use timeouts
	// WriteTimeout & ReadTimeout
	return &http.Server{
		Addr:    addr,
		Handler: h,
	}
}

// Start odpre socket in zacne strecti http server s podporo HTTP/2 in TLS
// TODO:
// 	- redircet iz HTTP na HTTPS (?)
//	- spremeni certifikat da se localhost.crt uporablja le pri development
func Start(s *http.Server) {
	// Zazene http server, ignoriramo error http.ErrServerClosed,
	// saj se prozi vedno pri ugasanju streznika
	go func() {
		// Preveri ali obstaja domain.crt
		if _, certErr := os.Stat("certs/localhost.crt"); certErr != nil {
			log.Fatal("Datoteka s SSL certifikatom ni bila najdena")
		}

		// Preveri ali obstaja domain.key
		if _, keyErr := os.Stat("certs/localhost.crt"); keyErr != nil {
			log.Fatal("Datoteka s SSL kljucem ni bila najdena")
		}

		if err := s.ListenAndServeTLS("certs/localhost.crt", "certs/localhost.key"); err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()
}
