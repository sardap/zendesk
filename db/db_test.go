package db_test

import (
	"os"
	"testing"

	"github.com/sardap/zendesk/db"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	ticketsFile, _ := os.Open("db_testdata/tickets.json")

	_, err := db.Create(ticketsFile)

	assert.NoError(t, err, "Error creating databse")
}

func TestGetTicket(t *testing.T) {

	ticketsFile, _ := os.Open("db_testdata/tickets.json")

	db, _ := db.Create(ticketsFile)

	_, err := db.GetTicket("436bf9b0-1147-4c0a-8439-6f79833bff5b")

	assert.NoError(t, err, "Error getting first ticket")
}
