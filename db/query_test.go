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
			validField:   "_id",
			validMatch:   "101",
			invalidMatch: "10",
		},
		{
			resource:     db.ResourceUser,
			validField:   "_id",
			validMatch:   "101",
			invalidMatch: "10",
		},
		{
			resource:     db.ResourceTicket,
			validField:   "_id",
			validMatch:   "101",
			invalidMatch: "10",
		},
	}

	for _, testCase := range testCases {
		// valid match
		fullMatch := db.FulLMatchCondition{
			Resource: testCase.resource,
			Field:    testCase.validField,
			Match:    testCase.validMatch,
		}

		_, err := fullMatch.Resolve(database)
		assert.NoErrorf(t, err,
			"value not found on resouce %s field %s match %s",
			testCase.resource, testCase.validField, testCase.validMatch,
		)

		// no match found
		fullMatch = db.FulLMatchCondition{
			Resource: testCase.resource,
			Field:    testCase.validField,
			Match:    testCase.invalidMatch,
		}

		_, err = fullMatch.Resolve(database)
		assert.ErrorIsf(t, err, db.ErrNotFound,
			"value found when should not exist on resouce %s field %s match %s",
			testCase.resource, testCase.validField, testCase.validMatch,
		)

		// no match found
		fullMatch = db.FulLMatchCondition{
			Resource: testCase.resource,
			Field:    testCase.invalidMatch,
			Match:    testCase.invalidMatch,
		}

		_, err = fullMatch.Resolve(database)
		assert.ErrorIsf(t, err, db.ErrFieldMissing,
			"field found when should not exist on resouce %s field %s",
			testCase.resource, testCase.validField,
		)
	}
}
