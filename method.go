package paynow

// PaymentMethod identifies the mobile money method used for an express-checkout
// (mobile) transaction. Pass one of these to Client.SendMobile.
type PaymentMethod string

const (
	// MethodEcocash is the EcoCash mobile money method.
	MethodEcocash PaymentMethod = "ecocash"

	// MethodOneMoney is the OneMoney mobile money method.
	MethodOneMoney PaymentMethod = "onemoney"

	// MethodInnbucks is the InnBucks mobile money method. InnBucks responses
	// include an authorization code that is surfaced on InitResponse.InnBucks.
	MethodInnbucks PaymentMethod = "innbucks"
)

// String returns the method as its wire value.
func (m PaymentMethod) String() string {
	return string(m)
}
