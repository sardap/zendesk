package ticket_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/sardap/zendesk/ticket"
	"github.com/sardap/zendesk/utility"
	"github.com/stretchr/testify/assert"
)

func TestTicketJsonParse(t *testing.T) {

	var result ticket.Ticket
	ticketsJson, _ := os.ReadFile("ticket_testdata/ticket.json")
	json.Unmarshal(ticketsJson, &result)

	// Check first value
	createdAt, _ := time.Parse(utility.ZendeskTimeFormat, "2016-04-28T11:19:34 -10:00")
	expectedTicket := ticket.Ticket{
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
