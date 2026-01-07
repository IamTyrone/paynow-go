package paynow_test

import (
	"testing"

	"github.com/IamTyrone/paynow-go/types"
)

func TestTransactionStatus_IsPaid(t *testing.T) {
	tests := []struct {
		status types.TransactionStatus
		want   bool
	}{
		{types.StatusPaid, true},
		{types.StatusCreated, false},
		{types.StatusSent, false},
		{types.StatusPending, false},
		{types.StatusCancelled, false},
		{types.StatusFailed, false},
		{types.StatusRefunded, false},
		{types.StatusAwaitingDelivery, false},
		{types.StatusDelivered, false},
		{types.TransactionStatus("Unknown"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsPaid(); got != tt.want {
				t.Errorf("TransactionStatus(%q).IsPaid() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestTransactionStatus_IsPending(t *testing.T) {
	tests := []struct {
		status types.TransactionStatus
		want   bool
	}{
		{types.StatusCreated, true},
		{types.StatusSent, true},
		{types.StatusPending, true},
		{types.StatusPaid, false},
		{types.StatusCancelled, false},
		{types.StatusFailed, false},
		{types.StatusRefunded, false},
		{types.StatusAwaitingDelivery, false},
		{types.StatusDelivered, false},
		{types.TransactionStatus("Unknown"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsPending(); got != tt.want {
				t.Errorf("TransactionStatus(%q).IsPending() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestTransactionStatus_IsFailed(t *testing.T) {
	tests := []struct {
		status types.TransactionStatus
		want   bool
	}{
		{types.StatusCancelled, true},
		{types.StatusFailed, true},
		{types.StatusCreated, false},
		{types.StatusSent, false},
		{types.StatusPending, false},
		{types.StatusPaid, false},
		{types.StatusRefunded, false},
		{types.StatusAwaitingDelivery, false},
		{types.StatusDelivered, false},
		{types.TransactionStatus("Unknown"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsFailed(); got != tt.want {
				t.Errorf("TransactionStatus(%q).IsFailed() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestTransactionStatus_Values(t *testing.T) {
	if types.StatusCreated != "Created" {
		t.Errorf("StatusCreated = %q, want %q", types.StatusCreated, "Created")
	}
	if types.StatusSent != "Sent" {
		t.Errorf("StatusSent = %q, want %q", types.StatusSent, "Sent")
	}
	if types.StatusPending != "Pending" {
		t.Errorf("StatusPending = %q, want %q", types.StatusPending, "Pending")
	}
	if types.StatusPaid != "Paid" {
		t.Errorf("StatusPaid = %q, want %q", types.StatusPaid, "Paid")
	}
	if types.StatusCancelled != "Cancelled" {
		t.Errorf("StatusCancelled = %q, want %q", types.StatusCancelled, "Cancelled")
	}
	if types.StatusFailed != "Failed" {
		t.Errorf("StatusFailed = %q, want %q", types.StatusFailed, "Failed")
	}
	if types.StatusRefunded != "Refunded" {
		t.Errorf("StatusRefunded = %q, want %q", types.StatusRefunded, "Refunded")
	}
}

func TestPaymentMethod_Ecocash(t *testing.T) {
	if types.MethodEcocash != "ecocash" {
		t.Errorf("MethodEcocash = %q, want %q", types.MethodEcocash, "ecocash")
	}
}

func TestPaymentMethod_StringConversion(t *testing.T) {
	method := types.MethodEcocash
	str := string(method)
	if str != "ecocash" {
		t.Errorf("string(MethodEcocash) = %q, want %q", str, "ecocash")
	}
}
