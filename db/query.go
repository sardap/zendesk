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

type Resource string

const (
	ResourceOrganization Resource = "organization"
	ResourceUser         Resource = "user"
	ResourceTicket       Resource = "ticket"
)

type Condition interface {
	Resolve(db *DB) (bool, error)
}

type Query struct {
	Target     Resource
	Conditions []Condition
}

type FulLMatchCondition struct {
	Resource Resource
	Field    string
	Match    string
}

func (f *FulLMatchCondition) Resolve(db *DB) (bool, error) {
	var resource interface{}

	switch f.Resource {
	case ResourceOrganization:
		resource = Organization{}
	case ResourceUser:
		resource = User{}
	case ResourceTicket:
		resource = Ticket{}
	default:
		return false, errors.Wrapf(ErrInvalidResouce, "%s", f.Resource)
	}

	reflectInfo := reflect.TypeOf(resource)
	for i := 0; i < reflectInfo.NumField(); i++ {
		memeber := reflectInfo.Field(i)
		tag := string(memeber.Tag)
		tag = strings.TrimPrefix(tag, `json:"`)
		tag = strings.TrimSuffix(tag, `"`)

		// Target filed found
		if strings.ToLower(f.Field) == tag {
			switch f.Resource {
			case ResourceOrganization:
				for _, val := range db.orgs {
					v := reflect.ValueOf(resource).Field(i).Interface()
					targetValue := fmt.Sprintf("%v", v)
					return targetValue == f.Match, nil
				}
			case ResourceUser:
				resource = User{}
			case ResourceTicket:
				resource = Ticket{}
			}

			break
		}
	}

	return false, errors.Wrapf(ErrFieldMissing, "unable to find %s in %s", f.Field, f.Resource)
}
