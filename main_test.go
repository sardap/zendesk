package main_test

import (
	"os"
	"testing"

	zendesk "github.com/sardap/zendesk"
	"github.com/sardap/zendesk/db"
	"github.com/stretchr/testify/assert"
)

func TestParseArgsValid(t *testing.T) {
	os.Create("testdata/orgs.json")
	defer os.Remove("testdata/orgs.json")
	os.Create("testdata/users.json")
	defer os.Remove("testdata/users.json")
	os.Create("testdata/tickets.json")
	defer os.Remove("testdata/tickets.json")

	// Set args via env vars
	os.Setenv("ORGS_FILE", "testdata/orgs.json")
	os.Setenv("USERS_FILE", "testdata/users.json")
	os.Setenv("TICKETS_FILE", "testdata/tickets.json")
	os.Setenv("QUERY", "user name test")

	args, err := zendesk.ParseFlags()
	assert.NoError(t, err)
	expectedArgs := zendesk.Args{
		OrganizationsFile: "testdata/orgs.json",
		UsersFile:         "testdata/users.json",
		TicketsFile:       "testdata/tickets.json",
		Query: db.Query{
			Conditions: []db.Condition{
				&db.FulLMatchCondition{
					Resource:  db.ResourceUser,
					Connector: db.ConnectorTypeUnion,
					Field:     "name",
					Match:     "test",
				},
			},
		},
	}
	assert.Equal(t, expectedArgs, args)
}
