package postgres_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rubinda/biolog/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Globalne spremenljivke za instance serviceov
var userServiceTest *postgres.UserService = &postgres.UserService{}
var speciesServiceTest *postgres.SpeciesService = &postgres.SpeciesService{}

func TestMain(m *testing.M) {
	// Prebere konfiguracijsko datoteko znotraj mape /config
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/go/src/github.com/rubinda/biolog/config")
	err := viper.ReadInConfig()
	if err != nil {
		// Parametri znotraj datoteke so kljucni za delovanje in se morajo prebrati ob zagonu
		log.Panic("There was a problem with the config file: ", err)
		os.Exit(1)
	}

	// Zagotovi testno podatkovno bazo
	if err := ensureTestDatabase(); err != nil {
		log.Fatal("Test database creation error: ", err)
		os.Exit(1)
	}

	var serviceErr error = nil
	userServiceTest, serviceErr = createUserService()
	if serviceErr != nil {
		log.Fatal("Can't create UserService: ", serviceErr)
		os.Exit(1)
	}
	speciesServiceTest, serviceErr = createSpeciesService()
	if serviceErr != nil {
		log.Fatal("Can't create SpeciesService: ", serviceErr)
		os.Exit(1)
	}
	// Zazeni teste
	runTests := m.Run()
	// Zapri povezave na podatkovno bazo
	// (!) Ignorira napake pri zapiranju virov
	userServiceTest.DB.Close()
	speciesServiceTest.DB.Close()

	// Odstrani testno podatkovno bazo
	/*if err := dropTestDatabase(); err != nil {
		log.Fatal("Test database drop error: ", err)
	}*/
	os.Exit(runTests)
}

// OpenDBConnection kreira povezavo na testno podatkovno povezavo
func OpenDBConnection() (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		viper.GetString("database.username"), viper.GetString("database.password"),
		viper.GetString("database.testdb"), viper.GetString("database.host"),
		viper.GetInt("database.port"), viper.GetString("database.sslmode"))
	db, err := sqlx.Open("postgres", connString)

	return db, err
}

// EnsureTestDatabase vzame dump od aktivne baze in ustvari novo testno podatkovno bazo
// (!) Deluje na macOS 10.13.3 z zsh, glej datoteko za navodila za uporabo.
// TODO preveri za druge OS (Windows, Ubuntu ...)
func ensureTestDatabase() error {
	macOSCmd := exec.Command("/bin/zsh", "../scripts/initTestDB-macos.zsh")
	if shellErr := macOSCmd.Run(); shellErr != nil {
		return shellErr
	}
	return nil
}

// DropTestDatabase pobrise testno podatkovno bazo
// (!) Deluje na macOS 10.13.3 z zsh
// TODO preveri za druge OS (Windows, Ubuntu ...)
func dropTestDatabase() error {
	cmdString := "dropdb -U biolog " + viper.GetString("database.testdb")
	dropDBCmd := exec.Command("/bin/zsh", "-c", cmdString)
	dropDBCmd.Stdin = strings.NewReader(viper.GetString("database.password"))
	if dropErr := dropDBCmd.Run(); dropErr != nil {
		return dropErr
	}
	return nil
}

// CreateUserService ustvari nov UserService s povezavo na bazo
func createUserService() (*postgres.UserService, error) {
	db, err := OpenDBConnection()
	if err != nil {
		return nil, err
	}
	s := &postgres.UserService{DB: db}

	return s, nil
}

// CreateSpeciesService ustvari nov SpeciesService s povezavo na bazo
func createSpeciesService() (*postgres.SpeciesService, error) {
	db, err := OpenDBConnection()
	if err != nil {
		return nil, err
	}
	s := &postgres.SpeciesService{DB: db}
	return s, nil
}
