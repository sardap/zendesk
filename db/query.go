package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrFieldMissing   error
	ErrInvalidResouce error
)

func init() {
	ErrFieldMissing = fmt.Errorf("field doesn't exist in that type")
	ErrInvalidResouce = fmt.Errorf("invalid resouce given")
}

type ResourceType string

const (
	ResourceOrganization ResourceType = "organization"
	ResourceUser         ResourceType = "user"
	ResourceTicket       ResourceType = "ticket"
)

// Converts json struct tag into the associated field
func getFieldName(obj interface{}, field string) string {
	// This is slower then a switch but requires no extra code When we add a new fied to a struct
	reflectInfo := reflect.TypeOf(obj)
	for i := 0; i < reflectInfo.NumField(); i++ {
		memeber := reflectInfo.Field(i)
		tag := string(memeber.Tag)
		tag = strings.TrimPrefix(tag, `json:"`)
		tag = strings.TrimSuffix(tag, `"`)

		// Target field found
		if strings.ToLower(field) == tag {
			return memeber.Name
		}
	}

	return ""
}

type Condition interface {
	Resolve(db *DB) (interface{}, error)
}

type Query struct {
	Conditions []Condition
}

type FulLMatchCondition struct {
	Resource ResourceType
	Field    string
	Match    string
}

func (f *FulLMatchCondition) Resolve(db *DB) (interface{}, error) {

	getValue := func(value interface{}, fieldName string) string {
		reflectValue := reflect.Indirect(reflect.ValueOf(value))
		return fmt.Sprintf("%v", reflectValue.FieldByName(fieldName).Interface())
	}

	switch f.Resource {
	case ResourceOrganization:
		fieldName := getFieldName(Organization{}, f.Field)
		if fieldName == "" {
			return nil, errors.Wrapf(ErrFieldMissing, "unable to find %s in %s", f.Field, f.Resource)
		}
		for _, val := range db.orgs {
			if getValue(val, fieldName) == f.Match {
				return val, nil
			}
		}
	case ResourceUser:
		fieldName := getFieldName(User{}, f.Field)
		if fieldName == "" {
			return nil, errors.Wrapf(ErrFieldMissing, "unable to find %s in %s", f.Field, f.Resource)
		}
		for _, val := range db.users {
			if getValue(val, fieldName) == f.Match {
				return val, nil
			}
		}
	case ResourceTicket:
		fieldName := getFieldName(Ticket{}, f.Field)
		if fieldName == "" {
			return nil, errors.Wrapf(ErrFieldMissing, "unable to find %s in %s", f.Field, f.Resource)
		}
		for _, val := range db.tickets {
			if getValue(val, fieldName) == f.Match {
				return val, nil
			}
		}
	default:
		return false, errors.Wrapf(ErrInvalidResouce, "%s", f.Resource)
	}

	return false, errors.Wrapf(ErrNotFound, "no matches for %s with %s in %s", f.Field, f.Match, f.Resource)
}
