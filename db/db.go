package db

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var (
	ErrNotFound          error
	ErrInvalidForeignKey error
)

func init() {
	ErrNotFound = fmt.Errorf("no entry found")
	ErrInvalidForeignKey = fmt.Errorf("invalid foreign key")
}

type DB struct {
	tickets map[string]*Ticket
	orgs    map[int]*Organization
	users   map[int]*User
}

func (d *DB) GetOrganization(id int) (*Organization, error) {
	result, ok := d.orgs[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%d", id)
	}

	return result, nil
}

func (d *DB) GetUser(id int) (*User, error) {
	result, ok := d.users[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%d", id)
	}

	return result, nil
}

func (d *DB) GetTicket(id string) (*Ticket, error) {
	result, ok := d.tickets[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%s", id)
	}

	return result, nil
}

func (d *DB) AddOrganization(toAdd Organization) error {
	if toAdd.users == nil {
		toAdd.users = make([]int, 0)
	}
	if toAdd.tickets == nil {
		toAdd.tickets = make([]string, 0)
	}
	d.orgs[toAdd.ID] = &toAdd

	return nil
}

func (d *DB) AddUser(toAdd User) error {
	if toAdd.assignee == nil {
		toAdd.assignee = make([]string, 0)
	}
	if toAdd.submitter == nil {
		toAdd.submitter = make([]string, 0)
	}

	d.users[toAdd.ID] = &toAdd
	// resolve foreign keys
	if org, err := d.GetOrganization(toAdd.OrganizationID); err == nil {
		org.users = append(org.users, toAdd.ID)
	}

	return nil
}

func (d *DB) AddTicket(toAdd Ticket) error {
	d.tickets[toAdd.ID] = &toAdd

	// resolve foreign keys
	if org, err := d.GetOrganization(toAdd.OrganizationID); err == nil {
		org.tickets = append(org.tickets, toAdd.ID)
	}

	if usr, err := d.GetUser(toAdd.SubmitterID); err == nil {
		usr.submitter = append(usr.submitter, toAdd.ID)
	}

	if usr, err := d.GetUser(toAdd.AssigneeID); err == nil {
		usr.assignee = append(usr.assignee, toAdd.ID)
	}

	return nil
}

func Create(orgsReader, usersReader, ticketsReader io.Reader) (*DB, error) {
	result := &DB{
		tickets: make(map[string]*Ticket),
		orgs:    make(map[int]*Organization),
		users:   make(map[int]*User),
	}

	// Organizations
	orgsJson, err := io.ReadAll(orgsReader)
	if err != nil {
		return nil, err
	}

	var orgs []Organization
	if err := json.Unmarshal(orgsJson, &orgs); err != nil {
		return nil, err
	}
	for _, org := range orgs {
		if err := result.AddOrganization(org); err != nil {
			return nil, err
		}
	}

	// Users
	usersJson, err := io.ReadAll(usersReader)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(usersJson, &users); err != nil {
		return nil, err
	}
	for _, user := range users {
		if err := result.AddUser(user); err != nil {
			return nil, err
		}
	}

	// Tickets
	ticketsJson, err := io.ReadAll(ticketsReader)
	if err != nil {
		return nil, err
	}

	var tickets []Ticket
	if err := json.Unmarshal(ticketsJson, &tickets); err != nil {
		return nil, err
	}
	for _, ticket := range tickets {
		if err := result.AddTicket(ticket); err != nil {
			return nil, err
		}
	}

	return result, nil
}
