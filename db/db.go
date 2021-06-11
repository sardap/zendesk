package db

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/sardap/zendesk/org"
	"github.com/sardap/zendesk/ticket"
)

var (
	ErrNotFound error
)

func init() {
	ErrNotFound = fmt.Errorf("no entry found")
}

type DB struct {
	tickets map[string]ticket.Ticket
	orgs    map[int]org.Organization
}

func Create(ticketReader io.Reader, orgsReader io.Reader) (*DB, error) {
	reuslt := &DB{
		tickets: make(map[string]ticket.Ticket),
		orgs:    make(map[int]org.Organization),
	}

	// Tickets
	ticketsJson, err := io.ReadAll(ticketReader)
	if err != nil {
		return nil, err
	}

	var tickets []ticket.Ticket
	if err := json.Unmarshal(ticketsJson, &tickets); err != nil {
		return nil, err
	}
	for _, ticket := range tickets {
		reuslt.tickets[ticket.ID] = ticket
	}

	// Organizations
	orgsJson, err := io.ReadAll(orgsReader)
	if err != nil {
		return nil, err
	}

	var orgs []org.Organization
	if err := json.Unmarshal(orgsJson, &orgs); err != nil {
		return nil, err
	}
	for _, org := range orgs {
		reuslt.orgs[org.ID] = org
	}

	return reuslt, nil
}

func (d *DB) GetTicket(id string) (*ticket.Ticket, error) {
	result, ok := d.tickets[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%s", id)
	}

	return &result, nil
}

func (d *DB) GetOrganization(id int) (*org.Organization, error) {
	result, ok := d.orgs[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%d", id)
	}

	return &result, nil
}
