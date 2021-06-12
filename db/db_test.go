package db_test

import (
	"bytes"
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

func createBlankDb() *db.DB {
	emptyJson := "[]"
	reuslt, _ := db.Create(
		bytes.NewBufferString(emptyJson),
		bytes.NewBufferString(emptyJson),
		bytes.NewBufferString(emptyJson),
	)
	return reuslt
}

func createLoadedDB() *db.DB {
	orgsFile, usersFile, ticketsFile := getFiles()
	defer orgsFile.Close()
	defer usersFile.Close()
	defer ticketsFile.Close()

	result, _ := db.Create(orgsFile, usersFile, ticketsFile)
	return result
}

func TestAddAndGetOrganization(t *testing.T) {
	database := createBlankDb()

	expected := db.Organization{
		ID: 200,
	}

	err := database.AddOrganization(expected)
	assert.NoError(t, err, "error adding org to DB")

	org, err := database.GetOrganization(200)
	assert.NoError(t, err, "error getting org from DB")
	assert.Equal(t, expected.ID, org.ID, "org gotten missmatch")
}

func TestAddAndGetUser(t *testing.T) {
	database := createBlankDb()

	expected := db.User{
		ID:             100,
		OrganizationID: 200,
	}

	err := database.AddUser(expected)
	assert.NoError(t, err, "error adding user")

	usr, err := database.GetUser(100)
	assert.NoError(t, err, "error getting user from DB")
	assert.Equal(t, expected.ID, usr.ID, "usr gotten missmatch")
}

func TestAddAndGetTicket(t *testing.T) {
	database := createBlankDb()

	expected := db.Ticket{
		ID:             "cool",
		OrganizationID: 200,
		AssigneeID:     1,
		SubmitterID:    2,
	}

	err := database.AddTicket(expected)
	assert.NoError(t, err, "error adding ticket")

	ticket, err := database.GetTicket("cool")
	assert.NoError(t, err, "error getting ticket from DB")
	assert.Equal(t, expected, *ticket, "ticket gotten missmatch")
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
	database := createLoadedDB()

	_, err := database.GetOrganization(101)
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetOrganization(99)
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}

func TestGetUser(t *testing.T) {
	database := createLoadedDB()

	_, err := database.GetUser(1)
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetUser(1000)
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}

func TestGetTicket(t *testing.T) {
	database := createLoadedDB()

	_, err := database.GetTicket("436bf9b0-1147-4c0a-8439-6f79833bff5b")
	assert.NoError(t, err, "Error getting first ticket")

	_, err = database.GetTicket("mr-garbage")
	assert.ErrorIs(t, err, db.ErrNotFound, "Missing notFoundm Error")
}
