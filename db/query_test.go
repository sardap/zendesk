package db_test

import (
	"testing"

	"github.com/sardap/zendesk/db"
)

func TestFullMatchCond(t *testing.T) {
	database := createLoadedDB()

	x := db.FulLMatchCondition{
		Resource: db.ResourceOrganization,
		Field:    "_id",
		Match:    "101",
	}

	x.Resolve(database)
}
