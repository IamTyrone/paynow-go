package paynow_test

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/IamTyrone/paynow-go"
)

const testKey = "3e9c8b12-integration-key"

func newTestClient(doer paynow.Doer) *paynow.Client {
	return paynow.New("12345", testKey,
		paynow.WithResultURL("https://example.com/result"),
		paynow.WithReturnURL("https://example.com/return"),
		paynow.WithHTTPClient(doer),
	)
}

func paidPayment() *paynow.Payment {
	return paynow.NewPayment("INV-1", "buyer@example.com").Add("Item", 10.00)
}

func TestSend_Success(t *testing.T) {
	doer := &mockDoer{response: signResponse(testKey,
		field{"status", "Ok"},
		field{"browserurl", "https://www.paynow.co.zw/payment/confirm/1"},
		field{"pollurl", "https://www.paynow.co.zw/interface/poll/1"},
	)}

	resp, err := newTestClient(doer).Send(context.Background(), paidPayment())
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if !resp.Success {
		t.Error("Success = false, want true")
	}
	if !resp.HasRedirect || resp.RedirectURL == "" {
		t.Error("expected a redirect URL for a web transaction")
	}
	if resp.PollURL == "" {
		t.Error("expected a poll URL")
	}
}

func TestSend_HitsWebEndpoint(t *testing.T) {
	doer := &mockDoer{response: signResponse(testKey, field{"status", "Ok"}, field{"pollurl", "p"})}
	_, _ = newTestClient(doer).Send(context.Background(), paidPayment())

	if !strings.Contains(doer.capturedURL, "initiatetransaction") {
		t.Errorf("Send() posted to %q, want the web endpoint", doer.capturedURL)
	}
	values, _ := url.ParseQuery(doer.capturedBody)
	if values.Get("id") != "12345" {
		t.Errorf("request id = %q, want 12345", values.Get("id"))
	}
	if values.Get("hash") == "" {
		t.Error("request should include a hash")
	}
}

func TestSendMobile_Success(t *testing.T) {
	doer := &mockDoer{response: signResponse(testKey,
		field{"status", "Ok"},
		field{"pollurl", "https://www.paynow.co.zw/interface/poll/1"},
		field{"instructions", "Dial *151#"},
	)}

	resp, err := newTestClient(doer).SendMobile(context.Background(), paidPayment(), "0771234567", paynow.MethodEcocash)
	if err != nil {
		t.Fatalf("SendMobile() error = %v", err)
	}
	if resp.Instructions == "" {
		t.Error("expected USSD instructions to be surfaced")
	}

	if !strings.Contains(doer.capturedURL, "remotetransaction") {
		t.Errorf("SendMobile() posted to %q, want the mobile endpoint", doer.capturedURL)
	}
	values, _ := url.ParseQuery(doer.capturedBody)
	if values.Get("method") != "ecocash" || values.Get("phone") != "0771234567" {
		t.Errorf("mobile fields not sent: method=%q phone=%q", values.Get("method"), values.Get("phone"))
	}
}

func TestSendMobile_InnBucks(t *testing.T) {
	doer := &mockDoer{response: signResponse(testKey,
		field{"status", "Ok"},
		field{"pollurl", "https://www.paynow.co.zw/interface/poll/1"},
		field{"authorizationcode", "ABC123"},
		field{"authorizationexpires", "2026-01-01 00:00:00"},
	)}

	resp, err := newTestClient(doer).SendMobile(context.Background(), paidPayment(), "0771234567", paynow.MethodInnbucks)
	if err != nil {
		t.Fatalf("SendMobile() error = %v", err)
	}
	if resp.InnBucks == nil {
		t.Fatal("expected InnBucks details to be populated")
	}
	if resp.InnBucks.AuthorizationCode != "ABC123" {
		t.Errorf("AuthorizationCode = %q, want ABC123", resp.InnBucks.AuthorizationCode)
	}
	if !strings.HasSuffix(resp.InnBucks.DeepLinkURL, "ABC123") {
		t.Errorf("DeepLinkURL = %q, want it to end with the code", resp.InnBucks.DeepLinkURL)
	}
	if !strings.HasSuffix(resp.InnBucks.QRCodeURL, "ABC123") {
		t.Errorf("QRCodeURL = %q, want it to end with the code", resp.InnBucks.QRCodeURL)
	}
}

func TestSend_ValidationErrors(t *testing.T) {
	client := newTestClient(&mockDoer{})

	if _, err := client.Send(context.Background(), nil); !errors.Is(err, paynow.ErrNoPayment) {
		t.Errorf("Send(nil) error = %v, want ErrNoPayment", err)
	}

	empty := paynow.NewPayment("INV", "buyer@example.com")
	if _, err := client.Send(context.Background(), empty); !errors.Is(err, paynow.ErrEmptyCart) {
		t.Errorf("Send(empty cart) error = %v, want ErrEmptyCart", err)
	}

	free := paynow.NewPayment("INV", "buyer@example.com").Add("Free", 0)
	if _, err := client.Send(context.Background(), free); !errors.Is(err, paynow.ErrNonPositiveTotal) {
		t.Errorf("Send(zero total) error = %v, want ErrNonPositiveTotal", err)
	}
}

func TestSendMobile_RequiresValidEmail(t *testing.T) {
	client := newTestClient(&mockDoer{})
	p := paynow.NewPayment("INV", "not-an-email").Add("Item", 10.00)

	if _, err := client.SendMobile(context.Background(), p, "0771234567", paynow.MethodEcocash); !errors.Is(err, paynow.ErrInvalidEmail) {
		t.Errorf("SendMobile() error = %v, want ErrInvalidEmail", err)
	}
}

func TestSend_PaynowError(t *testing.T) {
	doer := &mockDoer{response: "status=Error&error=" + url.QueryEscape("Invalid integration id")}

	resp, err := newTestClient(doer).Send(context.Background(), paidPayment())

	var apiErr *paynow.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Send() error = %v, want *paynow.APIError", err)
	}
	if resp == nil || resp.Success {
		t.Error("expected a populated, unsuccessful InitResponse")
	}
	if !strings.Contains(apiErr.Message, "Invalid integration id") {
		t.Errorf("APIError.Message = %q", apiErr.Message)
	}
}

func TestSend_TransportError(t *testing.T) {
	doer := &mockDoer{err: errors.New("network down")}
	if _, err := newTestClient(doer).Send(context.Background(), paidPayment()); err == nil {
		t.Error("expected a transport error")
	}
}

func TestSend_HashMismatch(t *testing.T) {
	doer := &mockDoer{response: "status=Ok&pollurl=p&hash=WRONGHASH"}
	if _, err := newTestClient(doer).Send(context.Background(), paidPayment()); !errors.Is(err, paynow.ErrHashMismatch) {
		t.Errorf("Send() error = %v, want ErrHashMismatch", err)
	}
}
