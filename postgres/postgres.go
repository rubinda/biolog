// Package postgres zdruzuje vse metode, ki se nanasajo na neposreden dostop do podatkovne baze.
// Vsak service ima svojo datoteko, kjer se nahajajo tudi funkcije povezane z njim.
package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Dodatek za PostgreSQL
)

// Open inicializira povezavo na podatkovno bazo PostgreSQL
func Open(user, password, dbname, host, sslmode string, port int) (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		user, password, dbname, host, port, sslmode)
	DB, error := sqlx.Open("postgres", connString) // Inicializiraj povezavo na bazo
	return DB, error
}

// Close zapre povezavo na podatkovno bazo in vraca morebitno napako pri zapiranju
func Close(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
