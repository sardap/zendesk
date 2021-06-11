package ticket_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/sardap/zendesk/ticket"
	"github.com/stretchr/testify/assert"
)

func TestTicketJsonParse(t *testing.T) {

	var tickets []ticket.Ticket
	ticketsJson, _ := os.ReadFile("ticket_testdata/tickets.json")
	json.Unmarshal(ticketsJson, &tickets)

	assert.NotNil(t, tickets, "tickets not parsed")

	assert.Equal(t, len(tickets), 200, "expected 200 entires for tickets")

	// Check first value
	expectedFirstTicket := ticket.Ticket{
		ID:             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
		URL:            "http://initech.zendesk.com/api/v2/tickets/436bf9b0-1147-4c0a-8439-6f79833bff5b.json",
		ExternalID:     "9210cdc9-4bee-485f-a078-35396cd74063",
		CreatedAt:      "2016-04-28T11:19:34 -10:00",
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

	assert.Equal(t, tickets[0], expectedFirstTicket, "first ticket not parsed correctly check Ticket json tags")

	// Check last value
	expectedLastTicket := ticket.Ticket{
		ID:             "50dfc8bc-31de-411e-92bf-a6d6b9dfa490",
		URL:            "http://initech.zendesk.com/api/v2/tickets/50dfc8bc-31de-411e-92bf-a6d6b9dfa490.json",
		ExternalID:     "8bc8bee7-2d98-4b69-b4a9-4f348ff41fa3",
		CreatedAt:      "2016-03-08T09:44:54 -11:00",
		Type:           "task",
		Subject:        "A Problem in South Africa",
		Description:    "Esse nisi occaecat pariatur veniam culpa dolore anim elit aliquip. Cupidatat mollit nulla consectetur ullamco tempor esse.",
		Priority:       "high",
		Status:         "hold",
		SubmitterID:    43,
		AssigneeID:     54,
		OrganizationID: 103,
		Tags: []string{
			"Georgia",
			"Tennessee",
			"Mississippi",
			"Marshall Islands",
		},
		HasIncidents: true,
		DueAt:        "2016-08-03T09:17:37 -10:00",
		Via:          "voice",
	}

	assert.Equal(t, tickets[len(tickets)-1], expectedLastTicket, "last ticket not parsed correctly check Ticket json tags")
}
