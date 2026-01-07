// Package types defines the core types used by the Paynow SDK.
package types

// PaymentMethod represents the available payment methods.
type PaymentMethod string

const (
	// MethodEcocash represents the EcoCash mobile money payment method.
	MethodEcocash PaymentMethod = "ecocash"
)
