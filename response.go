package paynow

import "strconv"

// InitResponse is the result of initiating a transaction with Client.Send or
// Client.SendMobile.
type InitResponse struct {
	// Status is the raw status Paynow returned (for example "Ok" or "Error").
	Status string

	// Success reports whether Paynow accepted the request.
	Success bool

	// HasRedirect reports whether RedirectURL is set. For web transactions the
	// customer must be redirected there to complete payment.
	HasRedirect bool

	// RedirectURL is the URL to send the customer to so they can pay. Only set
	// for web transactions.
	RedirectURL string

	// PollURL is the URL to poll (via Client.PollTransaction) to check the
	// transaction's status.
	PollURL string

	// Instructions holds USSD push instructions for the customer to dial, for
	// some mobile money payments.
	Instructions string

	// Error holds Paynow's error message when Success is false.
	Error string

	// Hash is the raw hash Paynow sent with the response.
	Hash string

	// InnBucks holds InnBucks-specific payment details when the response is for
	// an InnBucks transaction, and is nil otherwise.
	InnBucks *InnBucksInfo

	// Raw exposes every field Paynow returned, for access to fields the SDK does
	// not model explicitly.
	Raw map[string]string
}

// InnBucksInfo holds the details needed to complete an InnBucks payment. Paynow
// returns an authorization code which is turned into a tappable deep link and a
// scannable QR code.
type InnBucksInfo struct {
	// AuthorizationCode is the InnBucks payment code.
	AuthorizationCode string

	// DeepLinkURL opens the InnBucks app pre-filled with the payment code.
	DeepLinkURL string

	// QRCodeURL renders a QR code encoding the payment code.
	QRCodeURL string

	// ExpiresAt is when the authorization code expires, as returned by Paynow.
	ExpiresAt string
}

// newInitResponse builds an InitResponse from parsed, hash-verified values.
func newInitResponse(ov *orderedValues) *InitResponse {
	status, _ := ov.get("status")

	resp := &InitResponse{
		Status: status,
		Raw:    ov.asMap(),
	}
	resp.Success = !equalFoldTrim(status, responseError)

	if redirect, ok := ov.get("browserurl"); ok {
		resp.HasRedirect = true
		resp.RedirectURL = redirect
	}

	if !resp.Success {
		resp.Error, _ = ov.get("error")
		return resp
	}

	resp.PollURL, _ = ov.get("pollurl")
	resp.Hash, _ = ov.get("hash")
	if instructions, ok := ov.get("instructions"); ok {
		resp.Instructions = instructions
	}

	if code, ok := ov.get("authorizationcode"); ok && code != "" {
		expires, _ := ov.get("authorizationexpires")
		resp.InnBucks = &InnBucksInfo{
			AuthorizationCode: code,
			DeepLinkURL:       innbucksDeepLinkPrefix + code,
			QRCodeURL:         googleQRPrefix + code,
			ExpiresAt:         expires,
		}
	}

	return resp
}

// StatusResponse is the result of polling a transaction or processing a status
// update from Paynow.
type StatusResponse struct {
	// Status is the transaction's status as reported by Paynow.
	Status TransactionStatus

	// Paid is a convenience flag equivalent to Status.IsPaid.
	Paid bool

	// Amount is the transaction amount.
	Amount float64

	// Reference is the merchant's reference for the transaction.
	Reference string

	// PaynowReference is Paynow's own reference for the transaction.
	PaynowReference string

	// PollURL is the URL that can be polled to re-check the transaction status.
	PollURL string

	// Hash is the raw hash Paynow sent with the response.
	Hash string

	// Error holds Paynow's error message, if any.
	Error string

	// Raw exposes every field Paynow returned.
	Raw map[string]string
}

// newStatusResponse builds a StatusResponse from parsed values.
func newStatusResponse(ov *orderedValues) *StatusResponse {
	status, _ := ov.get("status")

	resp := &StatusResponse{
		Status: TransactionStatus(status),
		Raw:    ov.asMap(),
	}

	if equalFoldTrim(status, responseError) {
		resp.Error, _ = ov.get("error")
		return resp
	}

	resp.Paid = resp.Status.IsPaid()
	resp.Reference, _ = ov.get("reference")
	resp.PaynowReference, _ = ov.get("paynowreference")
	resp.PollURL, _ = ov.get("pollurl")
	resp.Hash, _ = ov.get("hash")

	if amount, ok := ov.get("amount"); ok {
		resp.Amount, _ = strconv.ParseFloat(amount, 64)
	}

	return resp
}
