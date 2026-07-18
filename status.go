package paynow

import "strings"

// TransactionStatus is the state of a transaction as reported by Paynow when
// polling a transaction or handling a status update. Comparisons are
// case-insensitive because Paynow is inconsistent about casing.
type TransactionStatus string

const (
	// StatusCreated indicates the transaction has been created but not yet sent.
	StatusCreated TransactionStatus = "Created"

	// StatusSent indicates the transaction has been sent to the payment provider.
	StatusSent TransactionStatus = "Sent"

	// StatusPending indicates the transaction is awaiting payment by the customer.
	StatusPending TransactionStatus = "Pending"

	// StatusPaid indicates the transaction has been paid successfully.
	StatusPaid TransactionStatus = "Paid"

	// StatusAwaitingDelivery indicates payment succeeded and delivery is awaited.
	StatusAwaitingDelivery TransactionStatus = "Awaiting Delivery"

	// StatusDelivered indicates the goods or services have been delivered.
	StatusDelivered TransactionStatus = "Delivered"

	// StatusCancelled indicates the transaction was cancelled.
	StatusCancelled TransactionStatus = "Cancelled"

	// StatusFailed indicates the transaction failed.
	StatusFailed TransactionStatus = "Failed"

	// StatusRefunded indicates the transaction was refunded.
	StatusRefunded TransactionStatus = "Refunded"

	// StatusDisputed indicates the transaction is under dispute.
	StatusDisputed TransactionStatus = "Disputed"
)

// Is reports whether s equals other, ignoring case.
func (s TransactionStatus) Is(other TransactionStatus) bool {
	return strings.EqualFold(string(s), string(other))
}

// IsPaid reports whether the transaction has been paid (or fulfilled beyond
// payment: awaiting delivery or delivered).
func (s TransactionStatus) IsPaid() bool {
	return s.Is(StatusPaid) || s.Is(StatusAwaitingDelivery) || s.Is(StatusDelivered)
}

// IsPending reports whether the transaction is still in progress and worth polling again.
func (s TransactionStatus) IsPending() bool {
	return s.Is(StatusCreated) || s.Is(StatusSent) || s.Is(StatusPending)
}

// IsFailed reports whether the transaction reached a terminal, unsuccessful state.
func (s TransactionStatus) IsFailed() bool {
	return s.Is(StatusCancelled) || s.Is(StatusFailed) || s.Is(StatusDisputed)
}
