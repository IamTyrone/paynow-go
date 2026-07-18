package paynow

import "context"

// Send initiates a normal web-based transaction. On success the returned
// InitResponse carries a RedirectURL the customer should be sent to in order to
// complete payment, and a PollURL for checking the transaction status.
//
// A non-nil error is returned for transport failures, malformed or
// hash-mismatched responses, and for business errors reported by Paynow (as an
// *APIError). Even when Paynow reports a business error the InitResponse is
// returned populated so its Error field can be inspected.
func (c *Client) Send(ctx context.Context, payment *Payment) (*InitResponse, error) {
	if err := validatePayment(payment); err != nil {
		return nil, err
	}

	body := c.buildWeb(payment).encode()
	return c.initiate(ctx, urlInitiateTransaction, body)
}

// SendMobile initiates an express-checkout mobile money transaction for the
// given phone number and method (for example paynow.MethodEcocash). Mobile
// transactions require a valid auth email on the payment.
//
// The returned InitResponse carries a PollURL and, depending on the method, USSD
// Instructions or InnBucks payment details. Error semantics match Send.
func (c *Client) SendMobile(ctx context.Context, payment *Payment, phone string, method PaymentMethod) (*InitResponse, error) {
	if err := validatePayment(payment); err != nil {
		return nil, err
	}
	if !isValidEmail(payment.AuthEmail) {
		return nil, ErrInvalidEmail
	}

	body := c.buildMobile(payment, phone, method).encode()
	return c.initiate(ctx, urlInitiateMobileTransaction, body)
}

// initiate posts a built request body to endpoint and parses the response into
// an InitResponse, verifying the hash on non-error responses.
func (c *Client) initiate(ctx context.Context, endpoint, body string) (*InitResponse, error) {
	raw, err := c.postForm(ctx, endpoint, body)
	if err != nil {
		return nil, err
	}

	values, err := parseResponse(raw)
	if err != nil {
		return nil, err
	}

	status, _ := values.get("status")
	if !equalFoldTrim(status, responseError) {
		if err := values.verifyHash(c.integrationKey); err != nil {
			return nil, err
		}
	}

	resp := newInitResponse(values)
	if !resp.Success {
		return resp, &APIError{Message: resp.Error}
	}
	return resp, nil
}

// validatePayment performs the sanity checks shared by Send and SendMobile.
func validatePayment(payment *Payment) error {
	if payment == nil {
		return ErrNoPayment
	}
	if len(payment.cart.items) == 0 {
		return ErrEmptyCart
	}
	if payment.Total() <= 0 {
		return ErrNonPositiveTotal
	}
	return nil
}
