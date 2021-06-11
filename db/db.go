package db

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
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
}

func Create(ticketReader io.ReadCloser) (*DB, error) {
	reuslt := &DB{
		tickets: make(map[string]ticket.Ticket),
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

	return reuslt, nil
}

func (d *DB) GetTicket(id string) (*ticket.Ticket, error) {
	result, ok := d.tickets[id]
	if !ok {
		return nil, errors.Wrapf(ErrNotFound, "%s", id)
	}

	return &result, nil
}
