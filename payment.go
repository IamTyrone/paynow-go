package paynow

// Payment is a transaction being built up before it is sent to Paynow. Create
// one with Client.CreatePayment, add one or more items with Add, then pass it to
// Client.Send or Client.SendMobile.
//
// A Payment behaves like a shopping cart: its total is the sum of the items
// added to it, and the item titles are sent to Paynow as additional info.
type Payment struct {
	// Reference is the merchant's unique identifier for the transaction.
	Reference string

	// AuthEmail is the customer's email address. It is required for mobile
	// (express checkout) transactions and optional for web transactions.
	AuthEmail string

	cart cart
}

// CreatePayment returns a new Payment with the given reference and auth email.
// The auth email may be empty for web transactions but is required for mobile ones.
func (c *Client) CreatePayment(reference, authEmail string) *Payment {
	return &Payment{Reference: reference, AuthEmail: authEmail}
}

// NewPayment returns a standalone Payment without needing a Client. It is
// equivalent to Client.CreatePayment and is convenient in tests or helpers.
func NewPayment(reference, authEmail string) *Payment {
	return &Payment{Reference: reference, AuthEmail: authEmail}
}

// Add appends an item to the payment's cart and returns the payment so calls
// can be chained. Quantity is optional and defaults to 1; only the first value
// is used if several are supplied.
func (p *Payment) Add(title string, amount float64, quantity ...int) *Payment {
	item := CartItem{Title: title, Amount: amount, Quantity: 1}
	if len(quantity) > 0 {
		item.Quantity = quantity[0]
	}
	p.cart.add(item)
	return p
}

// Items returns a copy of the items currently in the payment's cart.
func (p *Payment) Items() []CartItem {
	items := make([]CartItem, len(p.cart.items))
	copy(items, p.cart.items)
	return items
}

// Total returns the combined cost of every item in the cart.
func (p *Payment) Total() float64 {
	return p.cart.total()
}

// Info returns a comma-separated summary of the item titles in the cart.
func (p *Payment) Info() string {
	return p.cart.summary()
}
