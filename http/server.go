package http

import (
	"net/http"

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
// TODO redircet iz HTTP na HTTPS (?)
func Start(s *http.Server) {
	// Zazene http server, ignoriramo error http.ErrServerClosed,
	// saj se prozi vedno pri ugasanju streznika
	go func() {
		if err := s.ListenAndServeTLS("certs/domain.crt", "certs/domain.key"); err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()
}
