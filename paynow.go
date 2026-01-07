// Package paynow provides a Go SDK for integrating with Paynow Zimbabwe payment gateway.
//
// Example usage:
//
//	client := paynow.New(paynow.Config{
//		IntegrationID:  "your-integration-id",
//		IntegrationKey: "your-integration-key",
//		ResultURL:      "https://example.com/result",
//		ReturnURL:      "https://example.com/return",
//	})
//
//	response, err := client.SendMobile(paynow.Payment{
//		Reference: "INV-1001",
//		Amount:    10.00,
//		AuthEmail: "user@example.com",
//		Phone:     "0771234567",
//		Method:    paynow.MethodEcocash,
//	})
package paynow

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/IamTyrone/paynow-go/internal/hash"
	"github.com/IamTyrone/paynow-go/types"
)

const (
	defaultInitiateURL = "https://www.paynow.co.zw/interface/remotetransaction"
)

// HTTPClient interface allows for mocking HTTP requests in tests.
type HTTPClient interface {
	PostForm(url string, data url.Values) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

// Client represents a Paynow API client.
type Client struct {
	config     Config
	httpClient HTTPClient
}

// Config holds the configuration for the Paynow client.
type Config struct {
	IntegrationID  string
	IntegrationKey string
	ResultURL      string
	ReturnURL      string
}

// New creates a new Paynow client with the given configuration.
func New(config Config) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}
}

// NewWithHTTPClient creates a new Paynow client with a custom HTTP client.
// This is useful for testing or custom HTTP configurations.
func NewWithHTTPClient(config Config, httpClient HTTPClient) *Client {
	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

// Payment represents a payment request.
type Payment struct {
	Reference string
	Amount    float64
	AuthEmail string
	Phone     string
	Method    types.PaymentMethod
}

// InitResponse represents the response from initiating a payment.
type InitResponse struct {
	Status     string
	BrowserURL string
	PollURL    string
	Hash       string
	Error      string
}

// StatusResponse represents the response from polling payment status.
type StatusResponse struct {
	Reference       string
	Amount          float64
	PaynowReference string
	PollURL         string
	Status          types.TransactionStatus
	Hash            string
}

// SendMobile initiates a mobile money payment (e.g., EcoCash).
func (c *Client) SendMobile(payment Payment) (*InitResponse, error) {
	if err := c.validatePayment(payment); err != nil {
		return nil, err
	}

	if payment.Method == "" {
		payment.Method = types.MethodEcocash
	}

	data := c.buildRequestData(payment)

	requestHash := hash.GenerateHash(data, c.config.IntegrationKey)
	data.Set("hash", requestHash)

	response, err := c.httpClient.PostForm(defaultInitiateURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return c.parseInitResponse(string(body))
}

// PollTransaction checks the status of a transaction using the poll URL.
func (c *Client) PollTransaction(pollURL string) (*StatusResponse, error) {
	response, err := c.httpClient.Get(pollURL)
	if err != nil {
		return nil, fmt.Errorf("failed to poll transaction: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read poll response: %w", err)
	}

	return c.parseStatusResponse(string(body))
}

func (c *Client) validatePayment(payment Payment) error {
	if payment.Reference == "" {
		return fmt.Errorf("payment reference is required")
	}
	if payment.Amount <= 0 {
		return fmt.Errorf("payment amount must be greater than zero")
	}
	if payment.AuthEmail == "" {
		return fmt.Errorf("auth email is required")
	}
	if payment.Phone == "" {
		return fmt.Errorf("phone number is required for mobile payments")
	}
	return nil
}

func (c *Client) buildRequestData(payment Payment) url.Values {
	data := url.Values{}
	data.Set("id", c.config.IntegrationID)
	data.Set("reference", payment.Reference)
	data.Set("amount", fmt.Sprintf("%.2f", payment.Amount))
	data.Set("authemail", payment.AuthEmail)
	data.Set("phone", payment.Phone)
	data.Set("method", string(payment.Method))
	data.Set("returnurl", c.config.ReturnURL)
	data.Set("resulturl", c.config.ResultURL)
	data.Set("status", "Message")
	return data
}

func (c *Client) parseInitResponse(body string) (*InitResponse, error) {
	values, err := url.ParseQuery(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	status := values.Get("status")

	if status != "Error" {
		if err := hash.ValidateHash(body, c.config.IntegrationKey); err != nil {
			return nil, err
		}
	}

	resp := &InitResponse{
		Status:     status,
		BrowserURL: values.Get("browserurl"),
		PollURL:    values.Get("pollurl"),
		Hash:       values.Get("hash"),
		Error:      values.Get("error"),
	}

	if status == "Error" {
		return resp, fmt.Errorf("paynow error: %s", values.Get("error"))
	}

	return resp, nil
}

func (c *Client) parseStatusResponse(body string) (*StatusResponse, error) {
	values, err := url.ParseQuery(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	if err := hash.ValidateHash(body, c.config.IntegrationKey); err != nil {
		return nil, err
	}

	var amount float64
	if _, err := fmt.Sscanf(values.Get("amount"), "%f", &amount); err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	return &StatusResponse{
		Reference:       values.Get("reference"),
		Amount:          amount,
		PaynowReference: values.Get("paynowreference"),
		PollURL:         values.Get("pollurl"),
		Status:          types.TransactionStatus(values.Get("status")),
		Hash:            values.Get("hash"),
	}, nil
}
