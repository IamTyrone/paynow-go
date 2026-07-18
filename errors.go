package paynow

import (
	"errors"
	"fmt"
)

// Sentinel errors returned by the SDK. Callers can match against these with
// errors.Is to react to specific failure modes.
var (
	// ErrNoPayment is returned when a nil payment is passed to Send or SendMobile.
	ErrNoPayment = errors.New("paynow: payment is required")

	// ErrEmptyCart is returned when a payment has no items in its cart.
	ErrEmptyCart = errors.New("paynow: payment must contain at least one item")

	// ErrNonPositiveTotal is returned when a payment's total is not greater than zero.
	ErrNonPositiveTotal = errors.New("paynow: transaction total must be greater than zero")

	// ErrInvalidEmail is returned when a mobile payment is initiated without a
	// valid auth email. Mobile (express checkout) transactions require one.
	ErrInvalidEmail = errors.New("paynow: a valid auth email is required for mobile transactions")

	// ErrMissingHash is returned when a response from Paynow that should be
	// hashed does not contain a hash field.
	ErrMissingHash = errors.New("paynow: response does not contain a hash")

	// ErrHashMismatch is returned when the hash on a response from Paynow does
	// not match the hash computed locally, indicating a tampered or corrupt response.
	ErrHashMismatch = errors.New("paynow: response hash does not match")
)

// APIError represents a business error returned by Paynow itself, for example
// an invalid integration id or a malformed request. The Message field holds the
// human-readable reason supplied by Paynow.
type APIError struct {
	Message string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("paynow: %s", e.Message)
}
