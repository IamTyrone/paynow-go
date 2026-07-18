// Package paynow is a Go SDK for the Paynow Zimbabwe payment gateway.
//
// It supports the two ways of collecting a payment:
//
//   - Web transactions (Client.Send): the customer is redirected to Paynow to
//     choose how to pay.
//   - Mobile / express-checkout transactions (Client.SendMobile): the customer
//     pays directly with a mobile money method such as EcoCash or OneMoney.
//
// After initiating a payment you poll for its status with Client.PollTransaction,
// or handle the status update Paynow posts to your result URL with
// Client.ProcessStatusUpdate.
//
// Basic usage:
//
//	client := paynow.New(
//		"your-integration-id",
//		"your-integration-key",
//		paynow.WithResultURL("https://example.com/result"),
//		paynow.WithReturnURL("https://example.com/return"),
//	)
//
//	payment := client.CreatePayment("INV-1001", "customer@example.com")
//	payment.Add("T-shirt", 10.00)
//
//	resp, err := client.SendMobile(context.Background(), payment, "0771234567", paynow.MethodEcocash)
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(resp.PollURL)
package paynow

import "net/http"

// Doer is the subset of *http.Client the SDK needs. It lets callers inject a
// custom client (for timeouts, proxies, tracing or testing). *http.Client
// satisfies it out of the box.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client talks to the Paynow API on behalf of a single merchant integration.
// Create one with New. A Client is safe for concurrent use as long as the
// injected Doer is.
type Client struct {
	integrationID  string
	integrationKey string
	resultURL      string
	returnURL      string
	httpClient     Doer
}

// Option configures a Client. Pass options to New.
type Option func(*Client)

// WithResultURL sets the URL Paynow posts transaction status updates to.
func WithResultURL(url string) Option {
	return func(c *Client) { c.resultURL = url }
}

// WithReturnURL sets the URL the customer is returned to after paying.
func WithReturnURL(url string) Option {
	return func(c *Client) { c.returnURL = url }
}

// WithHTTPClient sets the HTTP client used for requests. By default a
// *http.Client with no timeout is used; supplying one with a timeout is
// recommended for production.
func WithHTTPClient(doer Doer) Option {
	return func(c *Client) {
		if doer != nil {
			c.httpClient = doer
		}
	}
}

// New creates a Client for the given integration credentials. Result and return
// URLs are optional here and can be supplied with WithResultURL / WithReturnURL
// or later with SetResultURL / SetReturnURL.
func New(integrationID, integrationKey string, opts ...Option) *Client {
	c := &Client{
		integrationID:  integrationID,
		integrationKey: integrationKey,
		httpClient:     &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// SetResultURL sets the URL Paynow posts transaction status updates to.
func (c *Client) SetResultURL(url string) { c.resultURL = url }

// SetReturnURL sets the URL the customer is returned to after paying.
func (c *Client) SetReturnURL(url string) { c.returnURL = url }
