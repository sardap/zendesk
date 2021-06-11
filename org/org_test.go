package org_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/sardap/zendesk/org"
	"github.com/sardap/zendesk/utility"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationJsonParse(t *testing.T) {

	var result org.Organization
	ticketsJson, _ := os.ReadFile("org_testdata/organization.json")
	json.Unmarshal(ticketsJson, &result)

	// Check first value
	createdAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-05-21T11:10:28 -10:00")
	expectedOrg := org.Organization{
		ID:         101,
		URL:        "http://initech.zendesk.com/api/v2/organizations/101.json",
		ExternalID: "9270ed79-35eb-4a38-a46f-35725197ea8d",
		Name:       "Enthaze",
		DomainNames: []string{
			"kage.com",
			"ecratic.com",
			"endipin.com",
			"zentix.com",
		},
		CreatedAt:     utility.ZendeskTime{Time: createdAt},
		Details:       "MegaCorp",
		SharedTickets: false,
		Tags: []string{
			"Fulton",
			"West",
			"Rodriguez",
			"Farley",
		},
	}

	assert.Equal(t, result, expectedOrg, "organization not parsed correctly check Organization json tags")
}
