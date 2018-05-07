package postgres_test

import (
	"testing"

	"github.com/rubinda/biolog"
	"github.com/stretchr/testify/assert"
)

var globalFalse bool = false

// TestCreateUsers preveri kreiranje uporabnikov v bazi
// Preveri naslednje scenarije:
// 	- kreira uporabnika, pri katerem je podan 'DisplayName' in 'PublicObservations'
//	- kreira uporabnika, pri katerem je podano le polje 'DisplayName' s sumniki
func TestCreateUser(t *testing.T) {
	n1, n2 := "River Tam", "Srečko Žališnik"
	cases := []struct {
		User biolog.User
	}{
		{biolog.User{DisplayName: &n1, PublicObservations: &globalFalse}},
		{biolog.User{DisplayName: &n2}},
	}
	for _, c := range cases {
		newUser, createErr := userServiceTest.CreateUser(c.User)
		if assert.NoError(t, createErr) {
			assert.Equal(t, c.User.DisplayName, newUser.DisplayName)
			assert.Equal(t, c.User.PublicObservations, newUser.PublicObservations)
		}
	}

}

// TestUser preveri pridobivanje uporabnikov iz baze
// Preveri naslednje scenarije:
//	- uspesno pridobiti uporabnika preko ID
func TestUser(t *testing.T) {
	cases := []struct {
		ID int
	}{
		{
			ID: 1,
		},
	}
	for _, c := range cases {
		getUser, getErr := userServiceTest.User(c.ID)
		if assert.NoError(t, getErr) {
			actualUser := biolog.User{}
			selectErr := userServiceTest.DB.Get(&actualUser,
				`SELECT * FROM biolog_user WHERE id = $1 LIMIT 1`, c.ID)
			if assert.NoError(t, selectErr) {
				assert.Equal(t, actualUser, getUser)
			}
		}
	}
}

// TestUsers preveri pridobivanje seznama uporabnikov iz baze
// Preveri naslednje scenarije:
// 	- pridobi vse uporabnike v bazi
func TestUsers(t *testing.T) {
	userList, err := userServiceTest.Users()
	if assert.NoError(t, err) {
		users := []*biolog.User{}
		selectErr := userServiceTest.DB.Select(&users, `SELECT * FROM biolog_user`)
		if assert.NoError(t, selectErr) {
			assert.EqualValues(t, users, userList)
		}
	}
}

// TestDeleteUser preveri brisanje uporabnika iz baze
// TODO preveri assert z us.DB.SELECT, pa primerjaj ce je prazno
// Preveri naslednje scenarije:
// 	- brise se uporabnik, ki obstaja
// 	- brise se uporabnik, ki ima zapise o opazanjih vrst
func TestDeleteUser(t *testing.T) {
	cases := []struct {
		ID         int
		ShouldStay bool
		Comment    string
	}{
		{
			ID:         2,
			ShouldStay: false,
			Comment:    "Can be deleted",
		},
		{
			ID:         3,
			ShouldStay: true,
			Comment:    "Can't delete, user has observation records",
		},
	}
	for _, c := range cases {
		userServiceTest.DeleteUser(c.ID)
		var userExists bool
		selectErr := userServiceTest.DB.QueryRow(`SELECT EXISTS
			(SELECT 1 FROM biolog_user WHERE id = $1 LIMIT 1)`, c.ID).Scan(&userExists)
		if assert.NoError(t, selectErr) {
			assert.Equal(t, c.ShouldStay, userExists)
		}
	}
}

// TestAuthProvider preveri pridobivanje podatkov o zunanjem ponudniku.
func TestAuthProvider(t *testing.T) {
	id := 1
	authProv, getErr := userServiceTest.AuthProvider(id)
	actualAuthProv := biolog.AuthProvider{}
	selectErr := userServiceTest.DB.Get(&actualAuthProv,
		`SELECT * FROM external_auth_provider WHERE id = $1 LIMIT 1`, id)
	if assert.NoError(t, getErr) {
		if assert.NoError(t, selectErr) {
			assert.Equal(t, actualAuthProv, authProv)
		}
	}
}
