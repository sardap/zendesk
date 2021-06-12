package db

import (
	"github.com/sardap/zendesk/utility"
)

type Organization struct {
	ID            int64               `json:"_id"`
	URL           string              `json:"url"`
	ExternalID    string              `json:"external_id"`
	Name          string              `json:"name"`
	DomainNames   []string            `json:"domain_names"`
	CreatedAt     utility.ZendeskTime `json:"created_at"`
	Details       string              `json:"details"`
	SharedTickets bool                `json:"shared_tickets"`
	Tags          []string            `json:"tags"`
	// foreign keys
	users   []int64
	tickets []string
}

func (o *Organization) GetUsers(db *DB) ([]*User, error) {
	var result []*User

	for _, id := range o.users {
		if usr, err := db.GetUser(id); err == nil {
			result = append(result, usr)
		}
	}

	return result, nil
}

func (o *Organization) GetTickets(db *DB) ([]*Ticket, error) {
	var result []*Ticket

	for _, id := range o.tickets {
		if ticket, err := db.GetTicket(id); err == nil {
			result = append(result, ticket)
		}
	}

	return result, nil
}

type User struct {
	ID             int64               `json:"_id"`
	URL            string              `json:"url"`
	ExternalID     string              `json:"external_id"`
	Name           string              `json:"name"`
	Alias          string              `json:"alias"`
	CreatedAt      utility.ZendeskTime `json:"created_at"`
	Active         bool                `json:"active"`
	Verified       bool                `json:"verified"`
	Shared         bool                `json:"shared"`
	Locale         string              `json:"locale"`
	Timezone       string              `json:"timezone"`
	LastLoginAt    utility.ZendeskTime `json:"last_login_at"`
	Email          string              `json:"email"`
	Phone          string              `json:"phone"`
	Signature      string              `json:"signature"`
	OrganizationID int64               `json:"organization_id"`
	Tags           []string            `json:"tags"`
	Suspended      bool                `json:"suspended"`
	Role           string              `json:"role"`
	// foreign keys
	assignee  []string
	submitter []string
}

func (u *User) GetAssignee(db *DB) ([]*Ticket, error) {
	var result []*Ticket

	for _, id := range u.assignee {
		if ticket, err := db.GetTicket(id); err == nil {
			result = append(result, ticket)
		}
	}

	return result, nil
}

func (u *User) GetSubmitter(db *DB) ([]*Ticket, error) {
	var result []*Ticket

	for _, id := range u.submitter {
		if ticket, err := db.GetTicket(id); err == nil {
			result = append(result, ticket)
		}
	}

	return result, nil
}

type Ticket struct {
	ID             string              `json:"_id"`
	URL            string              `json:"url"`
	ExternalID     string              `json:"external_id"`
	CreatedAt      utility.ZendeskTime `json:"created_at"`
	Type           string              `json:"type"`
	Subject        string              `json:"subject"`
	Description    string              `json:"description"`
	Priority       string              `json:"priority"`
	Status         string              `json:"status"`
	SubmitterID    int64               `json:"submitter_id"`
	AssigneeID     int64               `json:"assignee_id"`
	OrganizationID int64               `json:"organization_id"`
	Tags           []string            `json:"tags"`
	HasIncidents   bool                `json:"has_incidents"`
	DueAt          string              `json:"due_at"`
	Via            string              `json:"via"`
}
