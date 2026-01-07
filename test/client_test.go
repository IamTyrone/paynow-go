package paynow_test

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/IamTyrone/paynow-go"
	"github.com/IamTyrone/paynow-go/types"
)

// MockHTTPClient implements paynow.HTTPClient interface for testing.
type MockHTTPClient struct {
	PostFormFunc func(url string, data url.Values) (*http.Response, error)
	GetFunc      func(url string) (*http.Response, error)
}

func (m *MockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	if m.PostFormFunc != nil {
		return m.PostFormFunc(url, data)
	}
	return nil, errors.New("PostFormFunc not implemented")
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(url)
	}
	return nil, errors.New("GetFunc not implemented")
}

func newMockResponse(body string, statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func generateTestHash(values ...string) string {
	var hashString strings.Builder
	for _, v := range values {
		hashString.WriteString(v)
	}
	hash := sha512.Sum512([]byte(hashString.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func TestNew(t *testing.T) {
	config := paynow.Config{
		IntegrationID:  "test-id",
		IntegrationKey: "test-key",
		ResultURL:      "https://example.com/result",
		ReturnURL:      "https://example.com/return",
	}

	client := paynow.New(config)

	if client == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewWithHTTPClient(t *testing.T) {
	config := paynow.Config{
		IntegrationID:  "test-id",
		IntegrationKey: "test-key",
	}
	mockClient := &MockHTTPClient{}

	client := paynow.NewWithHTTPClient(config, mockClient)

	if client == nil {
		t.Fatal("NewWithHTTPClient() returned nil")
	}
}

func TestClient_SendMobile_ValidationErrors(t *testing.T) {
	client := paynow.New(paynow.Config{})

	tests := []struct {
		name    string
		payment paynow.Payment
		wantErr string
	}{
		{
			name:    "empty reference",
			payment: paynow.Payment{Amount: 10.00, AuthEmail: "test@example.com", Phone: "0771234567"},
			wantErr: "payment reference is required",
		},
		{
			name:    "zero amount",
			payment: paynow.Payment{Reference: "INV-001", Amount: 0, AuthEmail: "test@example.com", Phone: "0771234567"},
			wantErr: "payment amount must be greater than zero",
		},
		{
			name:    "negative amount",
			payment: paynow.Payment{Reference: "INV-001", Amount: -10.00, AuthEmail: "test@example.com", Phone: "0771234567"},
			wantErr: "payment amount must be greater than zero",
		},
		{
			name:    "empty email",
			payment: paynow.Payment{Reference: "INV-001", Amount: 10.00, Phone: "0771234567"},
			wantErr: "auth email is required",
		},
		{
			name:    "empty phone",
			payment: paynow.Payment{Reference: "INV-001", Amount: 10.00, AuthEmail: "test@example.com"},
			wantErr: "phone number is required for mobile payments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.SendMobile(tt.payment)
			if err == nil {
				t.Errorf("SendMobile() expected error containing %q, got nil", tt.wantErr)
			} else if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("SendMobile() error = %q, want error containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestClient_SendMobile_Success(t *testing.T) {
	integrationKey := "test-key-123"
	config := paynow.Config{
		IntegrationID:  "12345",
		IntegrationKey: integrationKey,
		ResultURL:      "https://example.com/result",
		ReturnURL:      "https://example.com/return",
	}

	browserURL := "https://www.paynow.co.zw/payment/confirm/123"
	pollURL := "https://www.paynow.co.zw/interface/poll/123"
	status := "Ok"

	hashInput := browserURL + pollURL + status + integrationKey
	responseHash := generateTestHash(hashInput)

	responseBody := fmt.Sprintf("browserurl=%s&pollurl=%s&status=%s&hash=%s",
		url.QueryEscape(browserURL),
		url.QueryEscape(pollURL),
		status,
		responseHash,
	)

	mockClient := &MockHTTPClient{
		PostFormFunc: func(reqURL string, data url.Values) (*http.Response, error) {
			if data.Get("id") != "12345" {
				t.Errorf("request id = %q, want %q", data.Get("id"), "12345")
			}
			if data.Get("hash") == "" {
				t.Error("request hash is empty")
			}
			return newMockResponse(responseBody, 200), nil
		},
	}

	client := paynow.NewWithHTTPClient(config, mockClient)

	response, err := client.SendMobile(paynow.Payment{
		Reference: "INV-001",
		Amount:    10.00,
		AuthEmail: "customer@example.com",
		Phone:     "0771234567",
		Method:    types.MethodEcocash,
	})

	if err != nil {
		t.Fatalf("SendMobile() unexpected error: %v", err)
	}
	if response.Status != "Ok" {
		t.Errorf("Status = %q, want %q", response.Status, "Ok")
	}
	if response.BrowserURL != browserURL {
		t.Errorf("BrowserURL = %q, want %q", response.BrowserURL, browserURL)
	}
	if response.PollURL != pollURL {
		t.Errorf("PollURL = %q, want %q", response.PollURL, pollURL)
	}
}

func TestClient_SendMobile_HTTPError(t *testing.T) {
	mockClient := &MockHTTPClient{
		PostFormFunc: func(url string, data url.Values) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	client := paynow.NewWithHTTPClient(paynow.Config{IntegrationKey: "key"}, mockClient)

	_, err := client.SendMobile(paynow.Payment{
		Reference: "INV-001",
		Amount:    10.00,
		AuthEmail: "test@example.com",
		Phone:     "0771234567",
	})

	if err == nil {
		t.Error("SendMobile() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to send request") {
		t.Errorf("error = %q, want error containing 'failed to send request'", err.Error())
	}
}

func TestClient_SendMobile_PaynowError(t *testing.T) {
	responseBody := "status=Error&error=Invalid+integration+id"

	mockClient := &MockHTTPClient{
		PostFormFunc: func(url string, data url.Values) (*http.Response, error) {
			return newMockResponse(responseBody, 200), nil
		},
	}

	client := paynow.NewWithHTTPClient(paynow.Config{IntegrationKey: "key"}, mockClient)

	response, err := client.SendMobile(paynow.Payment{
		Reference: "INV-001",
		Amount:    10.00,
		AuthEmail: "test@example.com",
		Phone:     "0771234567",
	})

	if err == nil {
		t.Error("SendMobile() expected error for Paynow error response")
	}
	if !strings.Contains(err.Error(), "paynow error") {
		t.Errorf("error = %q, want error containing 'paynow error'", err.Error())
	}
	if response == nil {
		t.Error("response should not be nil even on error")
	} else if response.Status != "Error" {
		t.Errorf("Status = %q, want %q", response.Status, "Error")
	}
}

func TestClient_SendMobile_DefaultsToEcocash(t *testing.T) {
	integrationKey := "key"
	responseHash := generateTestHash("Ok" + integrationKey)
	responseBody := "status=Ok&hash=" + responseHash

	var capturedMethod string
	mockClient := &MockHTTPClient{
		PostFormFunc: func(url string, data url.Values) (*http.Response, error) {
			capturedMethod = data.Get("method")
			return newMockResponse(responseBody, 200), nil
		},
	}

	client := paynow.NewWithHTTPClient(paynow.Config{IntegrationKey: integrationKey}, mockClient)

	_, _ = client.SendMobile(paynow.Payment{
		Reference: "INV-001",
		Amount:    10.00,
		AuthEmail: "test@example.com",
		Phone:     "0771234567",
	})

	if capturedMethod != "ecocash" {
		t.Errorf("method = %q, want %q (default)", capturedMethod, "ecocash")
	}
}

func TestClient_PollTransaction_Success(t *testing.T) {
	integrationKey := "test-key"
	reference := "INV-001"
	amount := "10.00"
	paynowRef := "12345678"
	pollURL := "https://www.paynow.co.zw/interface/poll/123"
	status := "Paid"

	hashInput := amount + paynowRef + pollURL + reference + status + integrationKey
	responseHash := generateTestHash(hashInput)

	responseBody := fmt.Sprintf(
		"amount=%s&paynowreference=%s&pollurl=%s&reference=%s&status=%s&hash=%s",
		amount, paynowRef, url.QueryEscape(pollURL), reference, status, responseHash,
	)

	mockClient := &MockHTTPClient{
		GetFunc: func(reqURL string) (*http.Response, error) {
			return newMockResponse(responseBody, 200), nil
		},
	}

	client := paynow.NewWithHTTPClient(paynow.Config{IntegrationKey: integrationKey}, mockClient)

	response, err := client.PollTransaction(pollURL)

	if err != nil {
		t.Fatalf("PollTransaction() unexpected error: %v", err)
	}
	if response.Status != types.StatusPaid {
		t.Errorf("Status = %q, want %q", response.Status, types.StatusPaid)
	}
	if response.Reference != reference {
		t.Errorf("Reference = %q, want %q", response.Reference, reference)
	}
	if response.Amount != 10.00 {
		t.Errorf("Amount = %v, want %v", response.Amount, 10.00)
	}
	if response.PaynowReference != paynowRef {
		t.Errorf("PaynowReference = %q, want %q", response.PaynowReference, paynowRef)
	}
}

func TestClient_PollTransaction_HTTPError(t *testing.T) {
	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, errors.New("connection refused")
		},
	}

	client := paynow.NewWithHTTPClient(paynow.Config{IntegrationKey: "key"}, mockClient)

	_, err := client.PollTransaction("https://example.com/poll")

	if err == nil {
		t.Error("PollTransaction() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to poll transaction") {
		t.Errorf("error = %q, want error containing 'failed to poll transaction'", err.Error())
	}
}

func TestClient_PollTransaction_InvalidHash(t *testing.T) {
	responseBody := "status=Paid&reference=INV-001&hash=INVALIDHASH"

	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return newMockResponse(responseBody, 200), nil
		},
	}

	client := paynow.NewWithHTTPClient(paynow.Config{IntegrationKey: "key"}, mockClient)

	_, err := client.PollTransaction("https://example.com/poll")

	if err == nil {
		t.Error("PollTransaction() expected hash validation error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid hash") {
		t.Errorf("error = %q, want error containing 'invalid hash'", err.Error())
	}
}

func TestPayment_Fields(t *testing.T) {
	payment := paynow.Payment{
		Reference: "INV-001",
		Amount:    99.99,
		AuthEmail: "test@example.com",
		Phone:     "0771234567",
		Method:    types.MethodEcocash,
	}

	if payment.Reference != "INV-001" {
		t.Errorf("Reference = %q, want %q", payment.Reference, "INV-001")
	}
	if payment.Amount != 99.99 {
		t.Errorf("Amount = %v, want %v", payment.Amount, 99.99)
	}
	if payment.AuthEmail != "test@example.com" {
		t.Errorf("AuthEmail = %q, want %q", payment.AuthEmail, "test@example.com")
	}
	if payment.Phone != "0771234567" {
		t.Errorf("Phone = %q, want %q", payment.Phone, "0771234567")
	}
	if payment.Method != types.MethodEcocash {
		t.Errorf("Method = %q, want %q", payment.Method, types.MethodEcocash)
	}
}

func TestConfig_Fields(t *testing.T) {
	config := paynow.Config{
		IntegrationID:  "id-123",
		IntegrationKey: "key-456",
		ResultURL:      "https://example.com/result",
		ReturnURL:      "https://example.com/return",
	}

	if config.IntegrationID != "id-123" {
		t.Errorf("IntegrationID = %q, want %q", config.IntegrationID, "id-123")
	}
	if config.IntegrationKey != "key-456" {
		t.Errorf("IntegrationKey = %q, want %q", config.IntegrationKey, "key-456")
	}
	if config.ResultURL != "https://example.com/result" {
		t.Errorf("ResultURL = %q, want %q", config.ResultURL, "https://example.com/result")
	}
	if config.ReturnURL != "https://example.com/return" {
		t.Errorf("ReturnURL = %q, want %q", config.ReturnURL, "https://example.com/return")
	}
}

func TestInitResponse_Fields(t *testing.T) {
	resp := paynow.InitResponse{
		Status:     "Ok",
		BrowserURL: "https://example.com/browser",
		PollURL:    "https://example.com/poll",
		Hash:       "ABC123",
		Error:      "",
	}

	if resp.Status != "Ok" {
		t.Errorf("Status = %q, want %q", resp.Status, "Ok")
	}
	if resp.BrowserURL != "https://example.com/browser" {
		t.Errorf("BrowserURL = %q, want %q", resp.BrowserURL, "https://example.com/browser")
	}
	if resp.PollURL != "https://example.com/poll" {
		t.Errorf("PollURL = %q, want %q", resp.PollURL, "https://example.com/poll")
	}
}

func TestStatusResponse_Fields(t *testing.T) {
	resp := paynow.StatusResponse{
		Reference:       "INV-001",
		Amount:          50.00,
		PaynowReference: "PN123",
		PollURL:         "https://example.com/poll",
		Status:          types.StatusPaid,
		Hash:            "HASH123",
	}

	if resp.Reference != "INV-001" {
		t.Errorf("Reference = %q, want %q", resp.Reference, "INV-001")
	}
	if resp.Amount != 50.00 {
		t.Errorf("Amount = %v, want %v", resp.Amount, 50.00)
	}
	if resp.PaynowReference != "PN123" {
		t.Errorf("PaynowReference = %q, want %q", resp.PaynowReference, "PN123")
	}
	if resp.Status != types.StatusPaid {
		t.Errorf("Status = %q, want %q", resp.Status, types.StatusPaid)
	}
}
