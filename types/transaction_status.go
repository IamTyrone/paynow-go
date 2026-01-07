package types

// TransactionStatus represents the status of a Paynow transaction.
type TransactionStatus string

const (
	// StatusCreated indicates the transaction has been created.
	StatusCreated TransactionStatus = "Created"

	// StatusSent indicates the transaction has been sent to the provider.
	StatusSent TransactionStatus = "Sent"

	// StatusPending indicates the transaction is awaiting payment.
	StatusPending TransactionStatus = "Pending"

	// StatusPaid indicates the transaction has been paid successfully.
	StatusPaid TransactionStatus = "Paid"

	// StatusCancelled indicates the transaction was cancelled.
	StatusCancelled TransactionStatus = "Cancelled"

	// StatusFailed indicates the transaction failed.
	StatusFailed TransactionStatus = "Failed"

	// StatusRefunded indicates the transaction was refunded.
	StatusRefunded TransactionStatus = "Refunded"

	// StatusAwaitingDelivery indicates awaiting delivery confirmation.
	StatusAwaitingDelivery TransactionStatus = "Awaiting Delivery"

	// StatusDelivered indicates the transaction has been delivered.
	StatusDelivered TransactionStatus = "Delivered"
)

// IsPaid returns true if the transaction status indicates successful payment.
func (s TransactionStatus) IsPaid() bool {
	return s == StatusPaid
}

// IsPending returns true if the transaction is still pending.
func (s TransactionStatus) IsPending() bool {
	return s == StatusCreated || s == StatusSent || s == StatusPending
}

// IsFailed returns true if the transaction failed or was cancelled.
func (s TransactionStatus) IsFailed() bool {
	return s == StatusCancelled || s == StatusFailed
}
