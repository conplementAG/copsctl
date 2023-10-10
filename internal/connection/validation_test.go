package connection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateConnectionString(t *testing.T) {
	data := []struct {
		testName           string
		connectionString   string
		environmentTag     string
		isTechnicalAccount bool
		shouldPanic        bool
	}{
		{
			testName:           "valid developer connection string",
			connectionString:   "https://example.com/env-developer.yaml",
			environmentTag:     "env",
			isTechnicalAccount: false,
			shouldPanic:        false,
		},
		{
			testName:           "invalid developer connection string",
			connectionString:   "https://example.com/env-developer.yaml",
			environmentTag:     "abc",
			isTechnicalAccount: false,
			shouldPanic:        true,
		},
		{
			testName:           "valid technical connection string",
			connectionString:   "https://example.com/env-technical-account.yaml",
			environmentTag:     "env",
			isTechnicalAccount: true,
			shouldPanic:        false,
		},
		{
			testName:           "invalid technical connection string",
			connectionString:   "https://example.com/env-technical-account.yaml",
			environmentTag:     "abc",
			isTechnicalAccount: true,
			shouldPanic:        true,
		},
		{
			testName:           "invalid developer connection string - invalid isTechnicalAccountFlag",
			connectionString:   "https://example.com/env-developer.yaml",
			environmentTag:     "env",
			isTechnicalAccount: true,
			shouldPanic:        true,
		},
		{
			testName:           "invalid technical connection string - invalid isTechnicalAccountFlag",
			connectionString:   "https://example.com/env-technical-account.yaml",
			environmentTag:     "env",
			isTechnicalAccount: false,
			shouldPanic:        true,
		},
	}

	for _, test := range data {
		t.Run(test.testName, func(t *testing.T) {
			// Arrange + Act + Assert
			if test.shouldPanic {
				assert.Panics(t, func() { validateConnectionString(test.connectionString, test.environmentTag, test.isTechnicalAccount) })

			} else {
				assert.NotPanics(t, func() { validateConnectionString(test.connectionString, test.environmentTag, test.isTechnicalAccount) })
			}
		})
	}
}
