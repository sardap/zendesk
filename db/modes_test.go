package db_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/sardap/zendesk/db"
	"github.com/sardap/zendesk/utility"
	"github.com/stretchr/testify/assert"
)

var (
	expectedOrg db.Organization
)

func init() {
	createdAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-05-21T11:10:28 -10:00")
	expectedOrg = db.Organization{
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

}

func TestOrganizationMatch(t *testing.T) {
	// URL
	match, err := expectedOrg.Match("url", "http://initech.zendesk.com/api/v2/organizations/101.json")
	assert.Truef(t, match, "should have matched url")
	assert.NoError(t, err, "error found for url")

	match, _ = expectedOrg.Match("url", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched url")

	// external_id
	match, err = expectedOrg.Match("external_id", "9270ed79-35eb-4a38-a46f-35725197ea8d")
	assert.Truef(t, match, "should have matched external_id")
	assert.NoError(t, err, "error found for external_id")

	match, _ = expectedOrg.Match("external_id", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched external_id")

	// name
	match, err = expectedOrg.Match("name", "Enthaze")
	assert.Truef(t, match, "should have matched name")
	assert.NoError(t, err, "error found for name")

	match, _ = expectedOrg.Match("name", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched name")

	// domain_names
	match, err = expectedOrg.Match("domain_names", "kage.com")
	assert.Truef(t, match, "should have matched domain_names")
	assert.NoError(t, err, "error found for domain_names")

	match, _ = expectedOrg.Match("domain_names", "sarda.dev")
	assert.Falsef(t, match, "should have not matched domain_names")

	// created_at
	match, err = expectedOrg.Match("created_at", "2016-05-21T11:10:28 -10:00")
	assert.Truef(t, match, "should have matched created_at")
	assert.NoError(t, err, "error found for created_at")

	match, _ = expectedOrg.Match("created_at", "2025-05-21T11:10:28 -10:00")
	assert.Falsef(t, match, "should have not matched created_at")

	_, err = expectedOrg.Match("created_at", "sarda.dev")
	assert.Error(t, err, "should have error for created_at")

	// details
	match, err = expectedOrg.Match("details", "MegaCorp")
	assert.Truef(t, match, "should have matched details")
	assert.NoError(t, err, "error found for details")

	match, _ = expectedOrg.Match("details", "sarda.dev")
	assert.Falsef(t, match, "should have not matched details")

	// shared_tickets
	match, err = expectedOrg.Match("shared_tickets", "false")
	assert.Truef(t, match, "should have matched shared_tickets")
	assert.NoError(t, err, "error found for shared_tickets")

	match, _ = expectedOrg.Match("shared_tickets", "true")
	assert.Falsef(t, match, "should have not matched shared_tickets")

	_, err = expectedOrg.Match("shared_tickets", "sarda.dev")
	assert.Error(t, err, "should have error for shared_tickets")

	// tags
	match, err = expectedOrg.Match("tags", "Fulton")
	assert.Truef(t, match, "should have matched tags")
	assert.NoError(t, err, "error found for tags")

	match, _ = expectedOrg.Match("tags", "sarda.dev")
	assert.Falsef(t, match, "should have not matched tags")
}

func TestOrganizationJsonParse(t *testing.T) {

	var result db.Organization
	ticketsJson, _ := os.ReadFile("db_testdata/organization.json")
	json.Unmarshal(ticketsJson, &result)

	assert.Equal(t, expectedOrg, result, "organization not parsed correctly check Organization json tags")
}

func TestUserJsonParse(t *testing.T) {

	var result db.User
	userJson, _ := os.ReadFile("db_testdata/user.json")
	json.Unmarshal(userJson, &result)

	// Check first value
	createdAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-04-15T05:19:46 -10:00")
	lastLoginAt, _ := time.Parse(utility.ZendeskTimeFormat, "2013-08-04T01:03:27 -10:00")
	expectedUser := db.User{
		ID:             1,
		URL:            "http://initech.zendesk.com/api/v2/users/1.json",
		ExternalID:     "74341f74-9c79-49d5-9611-87ef9b6eb75f",
		Name:           "Francisca Rasmussen",
		Alias:          "Miss Coffey",
		CreatedAt:      utility.ZendeskTime{Time: createdAt},
		Active:         true,
		Verified:       true,
		Shared:         false,
		Locale:         "en-AU",
		Timezone:       "Sri Lanka",
		LastLoginAt:    utility.ZendeskTime{Time: lastLoginAt},
		Email:          "coffeyrasmussen@flotonic.com",
		Phone:          "8335-422-718",
		Signature:      "Don't Worry Be Happy!",
		OrganizationID: 119,
		Tags: []string{
			"Springville",
			"Sutton",
			"Hartsville/Hartley",
			"Diaperville",
		},
		Suspended: true,
		Role:      "admin",
	}

	assert.Equal(t, expectedUser, result, "user not parsed correctly check User json tags")
}

func TestTicketJsonParse(t *testing.T) {

	var result db.Ticket
	ticketsJson, _ := os.ReadFile("db_testdata/ticket.json")
	json.Unmarshal(ticketsJson, &result)

	// Check first value
	createdAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-04-28T11:19:34 -10:00")
	expectedTicket := db.Ticket{
		ID:             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
		URL:            "http://initech.zendesk.com/api/v2/tickets/436bf9b0-1147-4c0a-8439-6f79833bff5b.json",
		ExternalID:     "9210cdc9-4bee-485f-a078-35396cd74063",
		CreatedAt:      utility.ZendeskTime{Time: createdAt},
		Type:           "incident",
		Subject:        "A Catastrophe in Korea (North)",
		Description:    "Nostrud ad sit velit cupidatat laboris ipsum nisi amet laboris ex exercitation amet et proident. Ipsum fugiat aute dolore tempor nostrud velit ipsum.",
		Priority:       "high",
		Status:         "pending",
		SubmitterID:    38,
		AssigneeID:     24,
		OrganizationID: 116,
		Tags: []string{
			"Ohio",
			"Pennsylvania",
			"American Samoa",
			"Northern Mariana Islands",
		},
		HasIncidents: false,
		DueAt:        "2016-07-31T02:37:50 -10:00",
		Via:          "web",
	}

	assert.Equal(t, result, expectedTicket, "ticket not parsed correctly check Ticket json tags")
}
