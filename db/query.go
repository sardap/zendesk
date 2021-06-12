package db

import (
	"fmt"

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

type Query struct {
	Conditions []Condition
}

func (q *Query) Resolve(db *DB) ([]Data, error) {
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

	var result []Data
	for _, val := range matches {
		result = append(result, val)
	}

	return result, nil
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
	case ResourceTicket:
	default:
		return nil, errors.Wrapf(ErrInvalidResouce, "%s", f.Resource)
	}

	return result, nil
}
