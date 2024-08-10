package requests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestToDomain_NoExpiration tests the ToDomain method without expiration
func TestToDomain_NoExpiration(t *testing.T) {
	request := CreateSecretRequest{
		SecretText:     "Test secret",
		ExpiresAfter:   0,
		RemainingViews: 5,
	}

	secret := request.ToDomain()

	assert.Equal(t, request.SecretText, secret.SecretText, "SecretText should match")
	assert.Equal(t, request.RemainingViews, secret.RemainingViews, "RemainingViews should match")

	// Check if the ExpiresAt is set to the end of the century
	expectedEndOfCentury := getEndOfCenturyDate()
	assert.Equal(t, expectedEndOfCentury.Year(), secret.ExpiresAt.Year(), "ExpiresAt year should match end of century")
	assert.Equal(t, expectedEndOfCentury.Month(), secret.ExpiresAt.Month(), "ExpiresAt month should match end of century")
	assert.Equal(t, expectedEndOfCentury.Day(), secret.ExpiresAt.Day(), "ExpiresAt day should match end of century")
}

// TestToDomain_WithExpiration tests the ToDomain method with expiration
func TestToDomain_WithExpiration(t *testing.T) {
	request := CreateSecretRequest{
		SecretText:     "Test secret",
		ExpiresAfter:   10, // 10 minutes
		RemainingViews: 5,
	}

	secret := request.ToDomain()

	assert.Equal(t, request.SecretText, secret.SecretText, "SecretText should match")
	assert.Equal(t, request.RemainingViews, secret.RemainingViews, "RemainingViews should match")

	// Check if the ExpiresAt is set correctly
	expectedExpiresAt := time.Now().UTC().Add(10 * time.Minute)
	assert.WithinDuration(t, expectedExpiresAt, secret.ExpiresAt, time.Minute, "ExpiresAt should be approximately 10 minutes from now")
}

// TestValidate_ValidRequest tests the Validate method for a valid request
func TestValidate_ValidRequest(t *testing.T) {
	request := CreateSecretRequest{
		SecretText:     "Valid secret",
		ExpiresAfter:   10,
		RemainingViews: 5,
	}

	err := request.Validate()

	assert.NoError(t, err, "Validate should not return an error for a valid request")
}

// TestValidate_InvalidRequests tests the Validate method for various invalid requests
func TestValidate_InvalidRequests(t *testing.T) {
	tests := []struct {
		name     string
		request  CreateSecretRequest
		expected string
	}{
		{
			name: "Empty Secret Text",
			request: CreateSecretRequest{
				SecretText:     "",
				ExpiresAfter:   10,
				RemainingViews: 5,
			},
			expected: "secret text is required",
		},
		{
			name: "Negative Expires After",
			request: CreateSecretRequest{
				SecretText:     "Negative expires",
				ExpiresAfter:   -5,
				RemainingViews: 5,
			},
			expected: "expires after should be greater than or equal to 0",
		},
		{
			name: "Negative Remaining Views",
			request: CreateSecretRequest{
				SecretText:     "Negative views",
				ExpiresAfter:   10,
				RemainingViews: -5,
			},
			expected: "remaining views should be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			assert.Error(t, err)
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

// TestGetEndOfCenturyDate tests the getEndOfCenturyDate function
func TestGetEndOfCenturyDate(t *testing.T) {
	endOfCentury := getEndOfCenturyDate()
	currentYear := time.Now().UTC().Year()

	assert.Equal(t, currentYear+100, endOfCentury.Year(), "The year should be 100 years from the current year")
	assert.Equal(t, time.December, endOfCentury.Month(), "The month should be December")
	assert.Equal(t, 31, endOfCentury.Day(), "The day should be the last day of December")
}
