package db

import "github.com/sardap/zendesk/utility"

type Organization struct {
	ID            int                 `json:"_id"`
	URL           string              `json:"url"`
	ExternalID    string              `json:"external_id"`
	Name          string              `json:"name"`
	DomainNames   []string            `json:"domain_names"`
	CreatedAt     utility.ZendeskTime `json:"created_at"`
	Details       string              `json:"details"`
	SharedTickets bool                `json:"shared_tickets"`
	Tags          []string            `json:"tags"`
}

type User struct {
	ID             int                 `json:"_id"`
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
	OrganizationID int                 `json:"organization_id"`
	Tags           []string            `json:"tags"`
	Suspended      bool                `json:"suspended"`
	Role           string              `json:"role"`
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
	SubmitterID    int                 `json:"submitter_id"`
	AssigneeID     int                 `json:"assignee_id"`
	OrganizationID int                 `json:"organization_id"`
	Tags           []string            `json:"tags"`
	HasIncidents   bool                `json:"has_incidents"`
	DueAt          string              `json:"due_at"`
	Via            string              `json:"via"`
}
