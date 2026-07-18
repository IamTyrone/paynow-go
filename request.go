package paynow

import "github.com/IamTyrone/paynow-go/internal/hash"

// buildWeb assembles the fields for a normal web-based transaction, in the order
// Paynow expects, and appends the request hash.
func (c *Client) buildWeb(payment *Payment) *orderedValues {
	data := newOrderedValues()
	data.set("resulturl", c.resultURL)
	data.set("returnurl", c.returnURL)
	data.set("reference", payment.Reference)
	data.set("amount", formatAmount(payment.Total()))
	data.set("id", c.integrationID)
	data.set("additionalinfo", payment.Info())
	data.set("authemail", payment.AuthEmail)
	data.set("status", "Message")

	c.sign(data)
	return data
}

// buildMobile assembles the fields for an express-checkout mobile transaction,
// in the order Paynow expects, and appends the request hash.
func (c *Client) buildMobile(payment *Payment, phone string, method PaymentMethod) *orderedValues {
	data := newOrderedValues()
	data.set("resulturl", c.resultURL)
	data.set("returnurl", c.returnURL)
	data.set("reference", payment.Reference)
	data.set("amount", formatAmount(payment.Total()))
	data.set("id", c.integrationID)
	data.set("additionalinfo", payment.Info())
	data.set("authemail", payment.AuthEmail)
	data.set("phone", phone)
	data.set("method", method.String())
	data.set("status", "Message")

	c.sign(data)
	return data
}

// sign computes the request hash over the current values and appends it.
func (c *Client) sign(data *orderedValues) {
	data.set("hash", hash.Make(data.signingValues(), c.integrationKey))
}
