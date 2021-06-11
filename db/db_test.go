package db_test

import (
	"os"
	"testing"

	"github.com/sardap/zendesk/db"
	"github.com/stretchr/testify/assert"
)

func getFiles() (*os.File, *os.File, *os.File) {
	organizationsFile, _ := os.Open("db_testdata/organizations.json")
	usersFile, _ := os.Open("db_testdata/users.json")
	ticketsFile, _ := os.Open("db_testdata/tickets.json")

	return organizationsFile, usersFile, ticketsFile
}

func TestCreate(t *testing.T) {
	orgsFile, usersFile, ticketsFile := getFiles()
	defer orgsFile.Close()
	defer usersFile.Close()
	defer ticketsFile.Close()

	_, err := db.Create(orgsFile, usersFile, ticketsFile)

	assert.NoError(t, err, "Error creating databse")
}

func TestGetOrganization(t *testing.T) {
	orgsFile, usersFile, ticketsFile := getFiles()
	defer orgsFile.Close()
	defer usersFile.Close()
	defer ticketsFile.Close()

	database, _ := db.Create(orgsFile, usersFile, ticketsFile)

	_, err := database.GetOrganization(101)
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetOrganization(99)
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}

func TestGetUser(t *testing.T) {
	orgsFile, usersFile, ticketsFile := getFiles()
	defer orgsFile.Close()
	defer usersFile.Close()
	defer ticketsFile.Close()

	database, _ := db.Create(orgsFile, usersFile, ticketsFile)

	_, err := database.GetUser(1)
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetUser(1000)
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}

func TestGetTicket(t *testing.T) {
	orgsFile, usersFile, ticketsFile := getFiles()
	defer orgsFile.Close()
	defer usersFile.Close()
	defer ticketsFile.Close()

	database, _ := db.Create(orgsFile, usersFile, ticketsFile)

	_, err := database.GetTicket("436bf9b0-1147-4c0a-8439-6f79833bff5b")
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetTicket("mr-garbage")
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}
