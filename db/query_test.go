package db_test

import (
	"testing"

	"github.com/sardap/zendesk/db"
	"github.com/stretchr/testify/assert"
)

func TestFulLMatchCondition(t *testing.T) {
	database := createLoadedDB()

	testCases := []struct {
		resource     db.ResourceType
		validField   string
		validMatch   string
		invalidMatch string
	}{
		{
			resource:     db.ResourceOrganization,
			validField:   "domain_names",
			validMatch:   "kage.com",
			invalidMatch: "sarda.dev",
		},
		{
			resource:     db.ResourceUser,
			validField:   "name",
			validMatch:   "Francisca Rasmussen",
			invalidMatch: "garbage garbagehead",
		},
		{
			resource:     db.ResourceTicket,
			validField:   "subject",
			validMatch:   "A Catastrophe in Korea (North)",
			invalidMatch: "garbage garbagehead",
		},
	}

	for _, testCase := range testCases {
		// valid match
		fullMatch := db.FulLMatchCondition{
			Resource: testCase.resource,
			Field:    testCase.validField,
			Match:    testCase.validMatch,
		}

		matches, err := fullMatch.Resolve(database)
		assert.NoErrorf(t, err,
			"value not found on resource %s field %s match %s",
			testCase.resource, testCase.validField, testCase.validMatch,
		)
		assert.GreaterOrEqual(t, 1, len(matches),
			"value not found on resource %s field %s match %s",
			testCase.resource, testCase.validField, testCase.validMatch,
		)

		// no match found
		fullMatch = db.FulLMatchCondition{
			Resource: testCase.resource,
			Field:    testCase.validField,
			Match:    testCase.invalidMatch,
		}

		matches, _ = fullMatch.Resolve(database)
		assert.Equalf(t, 0, len(matches),
			"value should not have been found resource %s field %s match %s",
			testCase.resource, testCase.validField, testCase.invalidMatch,
		)

		// field missing
		fullMatch = db.FulLMatchCondition{
			Resource: testCase.resource,
			Field:    testCase.invalidMatch,
			Match:    testCase.invalidMatch,
		}

		_, err = fullMatch.Resolve(database)
		assert.ErrorIsf(t, err, db.ErrFieldMissing,
			"field found when should not exist on resource %s field %s",
			testCase.resource, testCase.validField,
		)
	}
}

func TestQueryIntersection(t *testing.T) {
	database := createLoadedDB()

	query := db.Query{
		Conditions: []db.Condition{
			&db.FulLMatchCondition{
				Resource:  db.ResourceOrganization,
				Connector: db.ConnectorTypeIntersection,
				Field:     "details",
				Match:     "MegaCorp",
			},
			&db.FulLMatchCondition{
				Resource:  db.ResourceOrganization,
				Connector: db.ConnectorTypeIntersection,
				Field:     "domain_names",
				Match:     "otherway.com",
			},
		},
	}

	result, err := query.Resolve(database)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	query = db.Query{
		Conditions: []db.Condition{
			&db.FulLMatchCondition{
				Resource:  db.ResourceOrganization,
				Connector: db.ConnectorTypeIntersection,
				Field:     "details",
				Match:     "MegaCorp",
			},
		},
	}

	result, err = query.Resolve(database)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(result))
}

func TestQueryUnion(t *testing.T) {
	database := createLoadedDB()

	query := db.Query{
		Conditions: []db.Condition{
			&db.FulLMatchCondition{
				Resource:  db.ResourceOrganization,
				Connector: db.ConnectorTypeUnion,
				Field:     "details",
				Match:     "MegaCorp",
			},
			&db.FulLMatchCondition{
				Resource:  db.ResourceOrganization,
				Connector: db.ConnectorTypeUnion,
				Field:     "name",
				Match:     "Multron",
			},
		},
	}

	result, err := query.Resolve(database)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result))
}
