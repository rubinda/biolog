// Package postgres zdruzuje vse metode, ki se nanasajo na neposreden dostop do podatkovne baze.
// Vsak service ima svojo datoteko, kjer se nahajajo tudi funkcije povezane z njim.
package postgres

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Dodatek za PostgreSQL
)

// Konstanti za funkcijo, ki gradi SQL query
const buildUpdate string = "UPDATE"
const buildInsert string = "INSERT"

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

// CreateInsertQuery loops through the fields of an struct and buildz a INSERT INTO query
// Returns the query with bindvars and arguments for values
// Accepts all fields (if a pointer field is given it checks for non nil value)
// and ignores fields with the name id
func buildInsertUpdateQuery(queryType string, table string, o interface{}) (string, []interface{}) {

	// Query represents the actual query, and args are the values to add into the prepared statement
	var query strings.Builder
	var args []interface{}

	// Counts the arguments for bindvar puproses
	i := 1
	// Resets the commas in the first iteration of building
	firstIteration := true

	// Get the nonnil fields from the object into a map
	nnFields := getNonNilFields(o)

	switch queryType {
	case buildInsert:
		// Build the INSERT query based on non nil fields
		// queryVals are the bindvars for the values
		var queryVals strings.Builder
		fmt.Fprintf(&query, "INSERT INTO %s (", table)
		queryVals.WriteString("VALUES (")

		for key, value := range nnFields {
			// Do not prepend commas during the first iteration
			if firstIteration {
				fmt.Fprintf(&query, "%s", key)
				fmt.Fprintf(&queryVals, "$%d", i)
				firstIteration = false
			} else {
				fmt.Fprintf(&query, ", %s", key)
				fmt.Fprintf(&queryVals, ", $%d", i)
			}
			args = append(args, value)
			i++
		}
		queryVals.WriteString(") RETURNING *")
		fmt.Fprintf(&query, ") %s", queryVals.String())

	case buildUpdate:
		// Build the UPDATE query based on non nil fields
		fmt.Fprintf(&query, "UPDATE %s SET ", table)

		for key, value := range nnFields {
			// Don't include the ID in the update
			if key == "id" {
				continue
			}

			// Do not prepend the comma during the first iteration
			if firstIteration {
				fmt.Fprintf(&query, "%s = $%d", key, i)
				firstIteration = false
			} else {
				fmt.Fprintf(&query, ", %s = $%d", key, i)
			}
			args = append(args, value)
			i++
		}
		fmt.Fprintf(&query, " WHERE id = $%d", i)
		i++

	}
	return query.String(), args
}

// GetNonNilFields iterates over struct fields and returns
// lowercase field names and values
func getNonNilFields(o interface{}) map[string]interface{} {
	// Iterate over the fields with the help of reflect
	uTyp := reflect.TypeOf(o)
	uVal := reflect.ValueOf(o)
	fields := make(map[string]interface{})

	for i := 0; i < uTyp.NumField(); i++ {

		field := uTyp.Field(i)
		fieldVal := uVal.Field(i)

		// Check if the field is a pointer
		val := fieldVal.Interface()
		if fieldVal.Kind() == reflect.Ptr {

			// Skip the field if it is a pointer to nil,
			// otherwise get the correct value of it
			if fieldVal.IsNil() {
				continue
			} else {
				val = reflect.Indirect(fieldVal).Interface()
			}
		}

		// Get the struct field name
		fName := field.Tag.Get("db")
		// Fields can have tags (PascalCase vs snake_case)
		if fName == "" {
			fName = strings.ToLower(field.Name)
		}

		// Ignore the ID field
		if fName == "id" {
			continue
		}

		// Store the names and values of the fields
		fields[fName] = val
	}
	return fields
}
