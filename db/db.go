package db

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var (
	ErrNotFound error
)

func init() {
	ErrNotFound = fmt.Errorf("no entry found")
}

type DB struct {
	tickets map[string]Ticket
	orgs    map[int]Organization
	users   map[int]User
}

func Create(orgsReader, usersReader, ticketsReader io.Reader) (*DB, error) {
	reuslt := &DB{
		tickets: make(map[string]Ticket),
		orgs:    make(map[int]Organization),
		users:   make(map[int]User),
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
		reuslt.orgs[org.ID] = org
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
		reuslt.users[user.ID] = user
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
		reuslt.tickets[ticket.ID] = ticket
	}

	return reuslt, nil
}

func (d *DB) GetOrganization(id int) (*Organization, error) {
	result, ok := d.orgs[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%d", id)
	}

	return &result, nil
}

func (d *DB) GetUser(id int) (*User, error) {
	result, ok := d.users[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%d", id)
	}

	return &result, nil
}

func (d *DB) GetTicket(id string) (*Ticket, error) {
	result, ok := d.tickets[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%s", id)
	}

	return &result, nil
}
