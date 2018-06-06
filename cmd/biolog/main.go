// Package main Biolog API.
//
// Aplikacija nudi podporo pri popisu vrst po Sloveniji
//
//		Schemes: https
//		Host: localhost:4000
//		BasePath: /api/v1
// 		Version: 1.0.0
//		Contact: David Rubin<david.rubin95@gmail.com>
//
//		Consumes:
//		- application/json
//
//		Produces:
//		-application/json
//
// swagger:meta
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rubinda/biolog/http"
	"github.com/rubinda/biolog/postgres"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Force barve za logrus
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	// Prebere konfiguracijsko datoteko znotraj mape /config
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/go/src/github.com/rubinda/biolog/config")
	err := viper.ReadInConfig()
	if err != nil {
		// Parametri znotraj datoteke so kljucni za delovanje in se morajo prebrati ob zagonu
		log.Panic("There was a problem with the config file: ", err)
	}

	// Inicializira povezavo na podatkovno bazo s pomocjo konfiguracijske datoteke
	db, error := postgres.Open(viper.GetString("database.username"), viper.GetString("database.password"),
		viper.GetString("database.dbname"), viper.GetString("database.host"), viper.GetString("database.sslmode"),
		viper.GetInt("database.port"))
	if error != nil {
		log.Panic("Error while establishing database connection: ", error)
	}

	// Ustvari service in jim nastavi podatkovno povezavo
	us := &postgres.UserService{DB: db}
	ss := &postgres.SpeciesService{DB: db}
	// Dodaj instance service na handlerja
	h := http.NewRootHandler(us, ss)
	// Handlerju vpisi podatke za dostop do Google APIs
	h.OAuthConf.ClientID = viper.GetString("server.client-id")
	h.OAuthConf.ClientSecret = viper.GetString("server.client-secret")

	// Zazene nov streznik in caka na signal interrupt
	sAddr := ":" + viper.GetString("server.address")
	s := http.NewServer(sAddr, h)
	http.Start(s)
	log.Info("Server is running @ localhost:" + viper.GetString("server.address"))

	// Registrira poslusalca za signalom 'INTERRUPT' (Ctrl-C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs // Posreduj signal programu

	// Ugasni streznik z timeout
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	log.Info("Shutdown with timeout: ", timeout)
	defer cancel()
	if err = s.Shutdown(ctx); err != nil {
		log.Error(err)
	} else {
		log.Info("Server stopped")
	}

}
