package db_test

import (
	"os"
	"testing"

	"github.com/sardap/zendesk/db"
	"github.com/stretchr/testify/assert"
)

// returns ticketFileReader, OrgReader
func getFiles() (*os.File, *os.File) {
	ticketsFile, _ := os.Open("db_testdata/tickets.json")
	organizationsFile, _ := os.Open("db_testdata/organizations.json")

	return ticketsFile, organizationsFile
}

func TestCreate(t *testing.T) {
	ticketsFile, orgsFile := getFiles()
	defer ticketsFile.Close()
	defer orgsFile.Close()

	_, err := db.Create(ticketsFile, orgsFile)

	assert.NoError(t, err, "Error creating databse")
}

func TestGetTicket(t *testing.T) {
	ticketsFile, orgsFile := getFiles()
	defer ticketsFile.Close()
	defer orgsFile.Close()

	database, _ := db.Create(ticketsFile, orgsFile)

	_, err := database.GetTicket("436bf9b0-1147-4c0a-8439-6f79833bff5b")
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetTicket("mr-garbage")
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}

func TestGetOrganization(t *testing.T) {
	ticketsFile, orgsFile := getFiles()
	defer ticketsFile.Close()
	defer orgsFile.Close()

	database, _ := db.Create(ticketsFile, orgsFile)

	_, err := database.GetOrganization(101)
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetOrganization(99)
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}
