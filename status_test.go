package paynow_test

import (
	"testing"

	"github.com/IamTyrone/paynow-go"
)

func TestTransactionStatus_Predicates(t *testing.T) {
	tests := []struct {
		status                paynow.TransactionStatus
		paid, pending, failed bool
	}{
		{paynow.StatusCreated, false, true, false},
		{paynow.StatusSent, false, true, false},
		{paynow.StatusPending, false, true, false},
		{paynow.StatusPaid, true, false, false},
		{paynow.StatusAwaitingDelivery, true, false, false},
		{paynow.StatusDelivered, true, false, false},
		{paynow.StatusCancelled, false, false, true},
		{paynow.StatusFailed, false, false, true},
		{paynow.StatusDisputed, false, false, true},
		{paynow.TransactionStatus("Unknown"), false, false, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsPaid(); got != tt.paid {
				t.Errorf("IsPaid() = %v, want %v", got, tt.paid)
			}
			if got := tt.status.IsPending(); got != tt.pending {
				t.Errorf("IsPending() = %v, want %v", got, tt.pending)
			}
			if got := tt.status.IsFailed(); got != tt.failed {
				t.Errorf("IsFailed() = %v, want %v", got, tt.failed)
			}
		})
	}
}

func TestTransactionStatus_CaseInsensitive(t *testing.T) {
	if !paynow.TransactionStatus("paid").IsPaid() {
		t.Error(`"paid" should be recognised as paid regardless of case`)
	}
}
