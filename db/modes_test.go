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

func TestOrganizationJsonParse(t *testing.T) {

	var result db.Organization
	ticketsJson, _ := os.ReadFile("db_testdata/organization.json")
	json.Unmarshal(ticketsJson, &result)

	// Check first value
	createdAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-05-21T11:10:28 -10:00")
	expectedOrg := db.Organization{
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
