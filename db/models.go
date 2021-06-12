package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sardap/zendesk/utility"
)

func matchBool(expected bool, value string) (bool, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return b == expected, nil
}

func matchTime(expected time.Time, value string) (bool, error) {
	t, err := time.Parse(utility.ZendeskTimeFormat, value)
	if err != nil {
		return false, errors.Wrapf(err, "time should be in %s format", utility.ZendeskTimeFormat)
	}
	return t.Equal(expected), nil
}

func matchStringArray(ary []string, value string) (bool, error) {
	for _, name := range ary {
		if name == value {
			return true, nil
		}
	}
	return false, nil
}

func matchInt64(i int64, value string) (bool, error) {
	parsedInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false, errors.Wrapf(err, "given value should be base 10")
	}

	return i == parsedInt, nil
}

type Data interface {
	GetKey() string
}

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

func (o *Organization) GetKey() string {
	return fmt.Sprintf("%d", o.ID)
}

func (o *Organization) Match(field, value string) (bool, error) {
	switch field {
	case "url":
		return o.URL == value, nil
	case "external_id":
		return o.ExternalID == value, nil
	case "name":
		return o.Name == value, nil
	case "domain_names":
		return matchStringArray(o.DomainNames, value)
	case "created_at":
		return matchTime(o.CreatedAt.Time, value)
	case "details":
		return o.Details == value, nil
	case "shared_tickets":
		return matchBool(o.SharedTickets, value)
	case "tags":
		return matchStringArray(o.Tags, value)
	}

	return false, ErrFieldMissing
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

func (u *User) GetKey() string {
	return fmt.Sprintf("%d", u.ID)
}

func (u *User) Match(field, value string) (bool, error) {
	switch field {
	case "url":
		return u.URL == value, nil
	case "external_id":
		return u.ExternalID == value, nil
	case "name":
		return u.Name == value, nil
	case "alias":
		return u.Alias == value, nil
	case "created_at":
		return matchTime(u.CreatedAt.Time, value)
	case "active":
		return matchBool(u.Active, value)
	case "verified":
		return matchBool(u.Verified, value)
	case "shared":
		return matchBool(u.Shared, value)
	case "locale":
		return u.Locale == value, nil
	case "timezone":
		return u.Timezone == value, nil
	case "last_login_at":
		return matchTime(u.LastLoginAt.Time, value)
	case "email":
		return u.Email == value, nil
	case "phone":
		return u.Phone == value, nil
	case "signature":
		return u.Signature == value, nil
	case "organization_id":
		return matchInt64(u.OrganizationID, value)
	case "tags":
		return matchStringArray(u.Tags, value)
	case "suspended":
		return matchBool(u.Suspended, value)
	case "role":
		return u.Role == value, nil
	}

	return false, ErrFieldMissing
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
