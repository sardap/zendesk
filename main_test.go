package main_test

import (
	"os"
	"testing"

	"github.com/namsral/flag"
	zendesk "github.com/sardap/zendesk"
	"github.com/sardap/zendesk/db"
	"github.com/stretchr/testify/assert"
)

func TestParseArgsFullMatch(t *testing.T) {
	// This is to prevent flag parsed twice error
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	// Setup files
	os.Create("testdata/orgs.json")
	defer os.Remove("testdata/orgs.json")
	os.Create("testdata/users.json")
	defer os.Remove("testdata/users.json")
	os.Create("testdata/tickets.json")
	defer os.Remove("testdata/tickets.json")

	// Set args via env vars
	os.Setenv("ORGS_FILE", "testdata/orgs.json")
	defer os.Unsetenv("ORGS_FILE")
	os.Setenv("USERS_FILE", "testdata/users.json")
	defer os.Unsetenv("USERS_FILE")
	os.Setenv("TICKETS_FILE", "testdata/tickets.json")
	defer os.Unsetenv("TICKETS_FILE")
	os.Setenv("QUERY", "user name test")
	defer os.Unsetenv("QUERY")

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

func TestParseArgsIdMatch(t *testing.T) {
	// This is to prevent flag parsed twice error
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	// Setup files
	os.Create("testdata/orgs.json")
	defer os.Remove("testdata/orgs.json")
	os.Create("testdata/users.json")
	defer os.Remove("testdata/users.json")
	os.Create("testdata/tickets.json")
	defer os.Remove("testdata/tickets.json")

	// Set args via env vars
	os.Setenv("ORGS_FILE", "testdata/orgs.json")
	defer os.Unsetenv("ORGS_FILE")
	os.Setenv("USERS_FILE", "testdata/users.json")
	defer os.Unsetenv("USERS_FILE")
	os.Setenv("TICKETS_FILE", "testdata/tickets.json")
	defer os.Unsetenv("TICKETS_FILE")
	os.Setenv("QUERY", "user id 100")
	defer os.Unsetenv("QUERY")

	args, err := zendesk.ParseFlags()
	assert.NoError(t, err)
	expectedArgs := zendesk.Args{
		OrganizationsFile: "testdata/orgs.json",
		UsersFile:         "testdata/users.json",
		TicketsFile:       "testdata/tickets.json",
		Query: db.Query{
			Conditions: []db.Condition{
				&db.IDMatchCondition{
					Resource: db.ResourceUser,
					Target:   "100",
				},
			},
		},
	}
	assert.Equal(t, expectedArgs, args)
}
