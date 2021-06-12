package db

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrFieldMissing   error
	ErrInvalidResouce error
	ErrInvalidMatch   error
)

func init() {
	ErrFieldMissing = fmt.Errorf("field doesn't exist in that type")
	ErrInvalidResouce = fmt.Errorf("invalid resouce given")
	ErrInvalidMatch = fmt.Errorf("invalid match given")
}

type ResourceType string

const (
	ResourceOrganization ResourceType = "organization"
	ResourceUser         ResourceType = "user"
	ResourceTicket       ResourceType = "ticket"
)

type ConnectorType string

const (
	ConnectorTypeIntersection ConnectorType = "intersection"
	ConnectorTypeUnion        ConnectorType = "union"
)

type Condition interface {
	Resolve(db *DB) ([]Data, error)
	GetResource() ResourceType
	GetConnector() ConnectorType
}

type QueryResult struct {
	Target  []Data `json:"target"`
	Related struct {
		Orgs    []Data `json:"organizations"`
		Users   []Data `json:"users"`
		Tickets []Data `json:"tickets"`
	} `json:"related"`
}

type Query struct {
	Conditions []Condition
}

func (q *Query) Resolve(db *DB) (*QueryResult, error) {
	matches := make(map[string]Data)

	for i, con := range q.Conditions {
		condMatches, err := con.Resolve(db)
		if err != nil {
			return nil, err
		}

		switch con.GetConnector() {
		case ConnectorTypeIntersection:
			intersection := make(map[string]Data)
			for _, val := range condMatches {
				if i == 0 {
					intersection[val.GetKey()] = val
				} else {
					if _, ok := matches[val.GetKey()]; ok {
						intersection[val.GetKey()] = val
					}
				}
			}

			matches = intersection
		case ConnectorTypeUnion:
			for _, val := range condMatches {
				matches[val.GetKey()] = val
			}
		default:
			panic(fmt.Errorf("unimplemented connector type"))
		}
	}

	var result QueryResult
	for _, val := range matches {
		result.Target = append(result.Target, val)
		for _, related := range val.GetRelated(db) {
			switch related.GetResourceType() {
			case ResourceOrganization:
				result.Related.Orgs = append(result.Related.Orgs, val)
			case ResourceUser:
				result.Related.Users = append(result.Related.Users, val)
			case ResourceTicket:
				result.Related.Tickets = append(result.Related.Tickets, val)
			}
		}
	}

	return &result, nil
}

type IDMatchCondition struct {
	Resource ResourceType
	Target   string
}

func (i *IDMatchCondition) GetConnector() ConnectorType {
	return ConnectorTypeUnion
}

func (i *IDMatchCondition) GetResource() ResourceType {
	return i.Resource
}

func (i *IDMatchCondition) Resolve(db *DB) ([]Data, error) {
	switch i.Resource {
	case ResourceOrganization:
		id, err := strconv.ParseInt(i.Target, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "given target ID must be a number and base 10")
		}
		result, err := db.GetOrganization(id)
		if err != nil {
			return nil, err
		}
		return []Data{result}, nil
	case ResourceUser:
		id, err := strconv.ParseInt(i.Target, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "given target ID must be a number and base 10")
		}
		result, err := db.GetUser(id)
		if err != nil {
			return nil, err
		}
		return []Data{result}, nil
	case ResourceTicket:
		result, err := db.GetTicket(i.Target)
		if err != nil {
			return nil, err
		}
		return []Data{result}, nil
	default:
		return nil, errors.Wrapf(ErrInvalidResouce, "%s", i.Resource)
	}

}

type FulLMatchCondition struct {
	Resource  ResourceType
	Connector ConnectorType
	Field     string
	Match     string
}

func (f *FulLMatchCondition) GetConnector() ConnectorType {
	return f.Connector
}

func (f *FulLMatchCondition) GetResource() ResourceType {
	return f.Resource
}

func (f *FulLMatchCondition) Resolve(db *DB) ([]Data, error) {
	var result []Data

	switch f.Resource {
	case ResourceOrganization:
		for _, val := range db.orgs {
			match, err := val.Match(f.Field, f.Match)
			if err != nil {
				return nil, err
			}
			if match {
				result = append(result, val)
			}
		}
	case ResourceUser:
		for _, val := range db.users {
			match, err := val.Match(f.Field, f.Match)
			if err != nil {
				return nil, err
			}
			if match {
				result = append(result, val)
			}
		}
	case ResourceTicket:
		for _, val := range db.tickets {
			match, err := val.Match(f.Field, f.Match)
			if err != nil {
				return nil, err
			}
			if match {
				result = append(result, val)
			}
		}
	default:
		return nil, errors.Wrapf(ErrInvalidResouce, "%s", f.Resource)
	}

	return result, nil
}
