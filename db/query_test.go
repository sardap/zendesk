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

func TestIDMatchCondition(t *testing.T) {
	database := createLoadedDB()

	testCases := []struct {
		resource  db.ResourceType
		validID   string
		invalidID string
	}{
		{
			resource:  db.ResourceOrganization,
			validID:   "101",
			invalidID: "10",
		},
		{
			resource:  db.ResourceUser,
			validID:   "10",
			invalidID: "1000",
		},
		{
			resource:  db.ResourceTicket,
			validID:   "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			invalidID: "A Catastrophe in Korea (North)",
		},
	}

	for _, testCase := range testCases {
		// valid match
		idMatch := db.IDMatchCondition{
			Resource: testCase.resource,
			Target:   testCase.validID,
		}

		matches, err := idMatch.Resolve(database)
		assert.NoErrorf(t, err,
			"value not found on resource %s id %s",
			testCase.resource, testCase.validID,
		)
		assert.GreaterOrEqual(t, 1, len(matches),
			"value not found on resource %s id %s",
			testCase.resource, testCase.validID,
		)

		// no target found
		idMatch = db.IDMatchCondition{
			Resource: testCase.resource,
			Target:   testCase.invalidID,
		}

		matches, err = idMatch.Resolve(database)
		assert.ErrorIsf(t, err, db.ErrNotFound,
			"value should not have been found resource %s target %s",
			testCase.resource, testCase.invalidID,
		)
		assert.Equalf(t, 0, len(matches),
			"value should not have been found resource %s target %s",
			testCase.resource, testCase.invalidID,
		)
	}
}

func TestQueryRelated(t *testing.T) {
	database := createLoadedDB()

	query := db.Query{
		Conditions: []db.Condition{
			&db.IDMatchCondition{
				Resource: db.ResourceUser,
				Target:   "1",
			},
		},
	}

	result, err := query.Resolve(database)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Target))
	assert.Equal(t, 1, len(result.Related.Orgs))
	assert.Equal(t, 4, len(result.Related.Tickets))
	assert.Equal(t, 0, len(result.Related.Users))
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
	assert.Equal(t, 1, len(result.Target))

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
	assert.Equal(t, 9, len(result.Target))
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
	assert.Equal(t, 10, len(result.Target))
}
