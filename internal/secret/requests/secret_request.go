package requests

import (
	"errors"
	"time"

	"github.com/nalawade41/secret-server/internal/domain"
)

type CreateSecretRequest struct {
	SecretText     string `form:"secret" json:"secret"`
	ExpiresAfter   int    `form:"expireAfter" json:"expireAfter"`
	RemainingViews int    `form:"expireAfterViews" json:"expireAfterViews"`
}

type GetSecretRequest struct {
	Hash string `param:"hash"`
}

// ToDomain method to transform to Domain.Secret struct
func (c CreateSecretRequest) ToDomain() domain.Secret {
	expiresAtUtc := getEndOfCenturyDate()

	// set expires at to 100 years form now
	if c.ExpiresAfter != 0 {
		expiresAtUtc = time.Now().UTC().Add(time.Duration(c.ExpiresAfter) * time.Minute)
	}

	// Create expiresAt using current time and duration
	return domain.Secret{
		SecretText:     c.SecretText,
		ExpiresAt:      expiresAtUtc,
		RemainingViews: c.RemainingViews,
		CreatedAt:      time.Now().UTC(),
	}
}

// Validate method to validate the request
func (c CreateSecretRequest) Validate() error {
	if c.SecretText == "" {
		return errors.New("secret text is required")
	}

	if c.ExpiresAfter < 0 {
		return errors.New("expires after should be greater than or equal to 0")
	}

	if c.RemainingViews < 0 {
		return errors.New("remaining views should be greater than 0")
	}

	return nil
}

func getEndOfCenturyDate() time.Time {
	// get end of century year
	endOfCenturyYear := time.Now().UTC().Year() + 100

	// create the end of century date
	endOfCenturyDate := time.Date(endOfCenturyYear, time.December, 31, 23, 59, 59, 0, time.UTC)

	return endOfCenturyDate
}
