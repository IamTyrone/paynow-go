package paynow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/IamTyrone/paynow-go"
)

func paidStatusBody() string {
	return signResponse(testKey,
		field{"reference", "INV-1"},
		field{"amount", "10.00"},
		field{"paynowreference", "PN-987"},
		field{"pollurl", "https://www.paynow.co.zw/interface/poll/1"},
		field{"status", "Paid"},
	)
}

func TestPollTransaction_Success(t *testing.T) {
	doer := &mockDoer{response: paidStatusBody()}

	resp, err := newTestClient(doer).PollTransaction(context.Background(), "https://www.paynow.co.zw/interface/poll/1")
	if err != nil {
		t.Fatalf("PollTransaction() error = %v", err)
	}
	if !resp.Paid || !resp.Status.IsPaid() {
		t.Errorf("Status = %q, Paid = %v; want a paid transaction", resp.Status, resp.Paid)
	}
	if resp.Amount != 10.00 {
		t.Errorf("Amount = %v, want 10.00", resp.Amount)
	}
	if resp.Reference != "INV-1" || resp.PaynowReference != "PN-987" {
		t.Errorf("references = %q / %q", resp.Reference, resp.PaynowReference)
	}
}

func TestPollTransaction_HashMismatch(t *testing.T) {
	doer := &mockDoer{response: "status=Paid&reference=INV-1&hash=WRONG"}
	if _, err := newTestClient(doer).PollTransaction(context.Background(), "poll"); !errors.Is(err, paynow.ErrHashMismatch) {
		t.Errorf("PollTransaction() error = %v, want ErrHashMismatch", err)
	}
}

func TestPollTransaction_TransportError(t *testing.T) {
	doer := &mockDoer{err: errors.New("boom")}
	if _, err := newTestClient(doer).PollTransaction(context.Background(), "poll"); err == nil {
		t.Error("expected a transport error")
	}
}

func TestProcessStatusUpdate_Success(t *testing.T) {
	client := newTestClient(&mockDoer{})

	resp, err := client.ProcessStatusUpdate(paidStatusBody())
	if err != nil {
		t.Fatalf("ProcessStatusUpdate() error = %v", err)
	}
	if !resp.Paid {
		t.Error("expected the status update to report a paid transaction")
	}
}

func TestProcessStatusUpdate_HashMismatch(t *testing.T) {
	client := newTestClient(&mockDoer{})
	if _, err := client.ProcessStatusUpdate("status=Paid&reference=INV-1&hash=WRONG"); !errors.Is(err, paynow.ErrHashMismatch) {
		t.Errorf("ProcessStatusUpdate() error = %v, want ErrHashMismatch", err)
	}
}
