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
	expectedOrg    db.Organization
	expectedUser   db.User
	expectedTicket db.Ticket
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

	// Check first value
	createdAt, _ = time.Parse(utility.ZendeskTimeFormat, "2016-04-15T05:19:46 -10:00")
	lastLoginAt, _ := time.Parse(utility.ZendeskTimeFormat, "2013-08-04T01:03:27 -10:00")
	expectedUser = db.User{
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

	createdAt, _ = time.Parse(utility.ZendeskTimeFormat, "2016-04-28T11:19:34 -10:00")
	dueAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-07-31T02:37:50 -10:00")
	expectedTicket = db.Ticket{
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
		DueAt:        utility.ZendeskTime{Time: dueAt},
		Via:          "web",
	}
}

func TestOrganizationMatch(t *testing.T) {
	// Invalid filed
	_, err := expectedOrg.Match("garbage", "garbage")
	assert.ErrorIs(t, err, db.ErrFieldMissing, "invlaid field")

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

func TestUserMatch(t *testing.T) {
	// Invalid filed
	_, err := expectedUser.Match("garbage", "garbage")
	assert.ErrorIs(t, err, db.ErrFieldMissing, "invlaid field")

	// URL
	match, err := expectedUser.Match("url", "http://initech.zendesk.com/api/v2/users/1.json")
	assert.Truef(t, match, "should have matched url")
	assert.NoError(t, err, "error found for url")

	match, _ = expectedUser.Match("url", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched url")

	// external_id
	match, err = expectedUser.Match("external_id", "74341f74-9c79-49d5-9611-87ef9b6eb75f")
	assert.Truef(t, match, "should have matched external_id")
	assert.NoError(t, err, "error found for external_id")

	match, _ = expectedUser.Match("external_id", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched external_id")

	// name
	match, err = expectedUser.Match("name", "Francisca Rasmussen")
	assert.Truef(t, match, "should have matched name")
	assert.NoError(t, err, "error found for name")

	match, _ = expectedUser.Match("name", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched name")

	// alias
	match, err = expectedUser.Match("alias", "Miss Coffey")
	assert.Truef(t, match, "should have matched alias")
	assert.NoError(t, err, "error found for alias")

	match, _ = expectedUser.Match("alias", "sarda.dev")
	assert.Falsef(t, match, "should have not matched alias")

	// created_at
	match, err = expectedUser.Match("created_at", "2016-04-15T05:19:46 -10:00")
	assert.Truef(t, match, "should have matched created_at")
	assert.NoError(t, err, "error found for created_at")

	match, _ = expectedUser.Match("created_at", "2025-05-21T11:10:28 -10:00")
	assert.Falsef(t, match, "should have not matched created_at")

	_, err = expectedUser.Match("created_at", "sarda.dev")
	assert.Error(t, err, "should have error for created_at")

	// active
	match, err = expectedUser.Match("active", "true")
	assert.Truef(t, match, "should have matched active")
	assert.NoError(t, err, "error found for active")

	match, _ = expectedUser.Match("active", "false")
	assert.Falsef(t, match, "should have not matched active")

	_, err = expectedUser.Match("active", "garbage")
	assert.Error(t, err, "should error for active")

	// verified
	match, err = expectedUser.Match("verified", "true")
	assert.Truef(t, match, "should have matched verified")
	assert.NoError(t, err, "error found for verified")

	match, _ = expectedUser.Match("verified", "false")
	assert.Falsef(t, match, "should have not matched verified")

	// shared
	match, err = expectedUser.Match("shared", "false")
	assert.Truef(t, match, "should have matched shared")
	assert.NoError(t, err, "error found for shared")

	match, _ = expectedUser.Match("shared", "true")
	assert.Falsef(t, match, "should have not matched shared")

	// locale
	match, err = expectedUser.Match("locale", "en-AU")
	assert.Truef(t, match, "should have matched locale")
	assert.NoError(t, err, "error found for locale")

	match, _ = expectedUser.Match("locale", "garbage")
	assert.Falsef(t, match, "should have not matched locale")

	// timezone
	match, err = expectedUser.Match("timezone", "Sri Lanka")
	assert.Truef(t, match, "should have matched timezone")
	assert.NoError(t, err, "error found for timezone")

	match, _ = expectedUser.Match("timezone", "garbage")
	assert.Falsef(t, match, "should have not matched timezone")

	// last_login_at
	match, err = expectedUser.Match("last_login_at", "2013-08-04T01:03:27 -10:00")
	assert.Truef(t, match, "should have matched last_login_at")
	assert.NoError(t, err, "error found for last_login_at")

	match, err = expectedUser.Match("last_login_at", "2025-08-04T01:03:27 -10:00")
	assert.NoError(t, err, "error found for last_login_at")
	assert.Falsef(t, match, "should have not matched last_login_at")

	// email
	match, err = expectedUser.Match("email", "coffeyrasmussen@flotonic.com")
	assert.Truef(t, match, "should have matched email")
	assert.NoError(t, err, "error found for email")

	match, err = expectedUser.Match("email", "garbage@garbage.com")
	assert.NoError(t, err, "error found for email")
	assert.Falsef(t, match, "should have not matched email")

	// phone
	match, err = expectedUser.Match("phone", "8335-422-718")
	assert.Truef(t, match, "should have matched phone")
	assert.NoError(t, err, "error found for phone")

	match, err = expectedUser.Match("phone", "8335-422-999")
	assert.NoError(t, err, "error found for phone")
	assert.Falsef(t, match, "should have not matched phone")

	// signature
	match, err = expectedUser.Match("signature", "Don't Worry Be Happy!")
	assert.Truef(t, match, "should have matched signature")
	assert.NoError(t, err, "error found for signature")

	match, err = expectedUser.Match("signature", "garbage")
	assert.NoError(t, err, "error found for signature")
	assert.Falsef(t, match, "should have not matched signature")

	// organization_id
	match, err = expectedUser.Match("organization_id", "119")
	assert.Truef(t, match, "should have matched organization_id")
	assert.NoError(t, err, "error found for organization_id")

	match, err = expectedUser.Match("organization_id", "100")
	assert.NoError(t, err, "error found for organization_id")
	assert.Falsef(t, match, "should have not matched organization_id")

	_, err = expectedUser.Match("organization_id", "garbage")
	assert.Error(t, err, "error should be found for organization_id")

	// tags
	match, err = expectedUser.Match("tags", "Springville")
	assert.Truef(t, match, "should have matched tags")
	assert.NoError(t, err, "error found for tags")

	match, _ = expectedUser.Match("tags", "sarda.dev")
	assert.Falsef(t, match, "should have not matched tags")

	// shared
	match, err = expectedUser.Match("suspended", "true")
	assert.Truef(t, match, "should have matched suspended")
	assert.NoError(t, err, "error found for suspended")

	match, _ = expectedUser.Match("suspended", "false")
	assert.Falsef(t, match, "should have not matched suspended")

	// role
	match, err = expectedUser.Match("role", "admin")
	assert.Truef(t, match, "should have matched role")
	assert.NoError(t, err, "error found for role")

	match, err = expectedUser.Match("role", "grabge")
	assert.NoError(t, err, "error found for role")
	assert.Falsef(t, match, "should have not matched role")
}

func TestUserJsonParse(t *testing.T) {

	var result db.User
	userJson, _ := os.ReadFile("db_testdata/user.json")
	json.Unmarshal(userJson, &result)

	assert.Equal(t, expectedUser, result, "user not parsed correctly check User json tags")
}

func TestTicketMatch(t *testing.T) {
	// Invalid filed
	_, err := expectedTicket.Match("garbage", "garbage")
	assert.ErrorIs(t, err, db.ErrFieldMissing, "invlaid field")

	// URL
	match, err := expectedTicket.Match("url", "http://initech.zendesk.com/api/v2/tickets/436bf9b0-1147-4c0a-8439-6f79833bff5b.json")
	assert.Truef(t, match, "should have matched url")
	assert.NoError(t, err, "error found for url")

	match, _ = expectedTicket.Match("url", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched url")

	// external_id
	match, err = expectedTicket.Match("external_id", "9210cdc9-4bee-485f-a078-35396cd74063")
	assert.Truef(t, match, "should have matched external_id")
	assert.NoError(t, err, "error found for external_id")

	match, _ = expectedTicket.Match("external_id", "https://sarda.dev")
	assert.Falsef(t, match, "should have not matched external_id")

	// created_at
	match, err = expectedTicket.Match("created_at", "2016-04-28T11:19:34 -10:00")
	assert.Truef(t, match, "should have matched created_at")
	assert.NoError(t, err, "error found for created_at")

	match, _ = expectedTicket.Match("created_at", "2025-05-21T11:10:28 -10:00")
	assert.Falsef(t, match, "should have not matched created_at")

	// type
	match, err = expectedTicket.Match("type", "incident")
	assert.Truef(t, match, "should have matched type")
	assert.NoError(t, err, "error found for type")

	match, _ = expectedTicket.Match("type", "false")
	assert.Falsef(t, match, "should have not matched type")

	// subject
	match, err = expectedTicket.Match("subject", "A Catastrophe in Korea (North)")
	assert.Truef(t, match, "should have matched subject")
	assert.NoError(t, err, "error found for subject")

	match, _ = expectedTicket.Match("subject", "false")
	assert.Falsef(t, match, "should have not matched subject")

	// description
	match, err = expectedTicket.Match("description", "Nostrud ad sit velit cupidatat laboris ipsum nisi amet laboris ex exercitation amet et proident. Ipsum fugiat aute dolore tempor nostrud velit ipsum.")
	assert.Truef(t, match, "should have matched description")
	assert.NoError(t, err, "error found for description")

	match, _ = expectedTicket.Match("description", "false")
	assert.Falsef(t, match, "should have not matched description")

	// priority
	match, err = expectedTicket.Match("priority", "high")
	assert.Truef(t, match, "should have matched priority")
	assert.NoError(t, err, "error found for priority")

	match, _ = expectedTicket.Match("priority", "false")
	assert.Falsef(t, match, "should have not matched priority")

	// status
	match, err = expectedTicket.Match("status", "pending")
	assert.Truef(t, match, "should have matched status")
	assert.NoError(t, err, "error found for status")

	match, _ = expectedTicket.Match("status", "false")
	assert.Falsef(t, match, "should have not matched status")

	// submitter_id
	match, err = expectedTicket.Match("submitter_id", "38")
	assert.Truef(t, match, "should have matched submitter_id")
	assert.NoError(t, err, "error found for submitter_id")

	match, _ = expectedTicket.Match("submitter_id", "10")
	assert.Falsef(t, match, "should have not matched submitter_id")

	// assignee_id
	match, err = expectedTicket.Match("assignee_id", "24")
	assert.Truef(t, match, "should have matched assignee_id")
	assert.NoError(t, err, "error found for assignee_id")

	match, _ = expectedTicket.Match("assignee_id", "10")
	assert.Falsef(t, match, "should have not matched assignee_id")

	// organization_id
	match, err = expectedTicket.Match("organization_id", "116")
	assert.Truef(t, match, "should have matched organization_id")
	assert.NoError(t, err, "error found for organization_id")

	match, _ = expectedTicket.Match("organization_id", "10")
	assert.Falsef(t, match, "should have not matched organization_id")

	// organization_id
	match, err = expectedTicket.Match("organization_id", "116")
	assert.Truef(t, match, "should have matched organization_id")
	assert.NoError(t, err, "error found for organization_id")

	match, _ = expectedTicket.Match("organization_id", "10")
	assert.Falsef(t, match, "should have not matched organization_id")

	// tags
	match, err = expectedTicket.Match("tags", "Pennsylvania")
	assert.Truef(t, match, "should have matched tags")
	assert.NoError(t, err, "error found for tags")

	match, _ = expectedTicket.Match("tags", "sarda.dev")
	assert.Falsef(t, match, "should have not matched tags")

	// has_incidents
	match, err = expectedTicket.Match("has_incidents", "false")
	assert.Truef(t, match, "should have matched has_incidents")
	assert.NoError(t, err, "error found for has_incidents")

	match, _ = expectedTicket.Match("has_incidents", "true")
	assert.Falsef(t, match, "should have not matched has_incidents")

	// has_incidents
	match, err = expectedTicket.Match("due_at", "2016-07-31T02:37:50 -10:00")
	assert.Truef(t, match, "should have matched due_at")
	assert.NoError(t, err, "error found for due_at")

	match, err = expectedTicket.Match("due_at", "2025-07-31T02:37:50 -10:00")
	assert.Falsef(t, match, "should have not matched due_at")
	assert.NoError(t, err, "error found for due_at")

	// via
	match, err = expectedTicket.Match("via", "web")
	assert.Truef(t, match, "should have matched via")
	assert.NoError(t, err, "error found for via")

	match, _ = expectedTicket.Match("via", "garbage")
	assert.Falsef(t, match, "should have not matched via")
}

func TestTicketJsonParse(t *testing.T) {

	var result db.Ticket
	ticketsJson, _ := os.ReadFile("db_testdata/ticket.json")
	json.Unmarshal(ticketsJson, &result)

	assert.Equal(t, result, expectedTicket, "ticket not parsed correctly check Ticket json tags")
}
