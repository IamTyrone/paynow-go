package paynow

import "strings"

// CartItem is a single line item in a Payment's cart.
type CartItem struct {
	// Title is the human-readable name of the item.
	Title string

	// Amount is the price of a single unit of the item.
	Amount float64

	// Quantity is the number of units. A value of zero is treated as one.
	Quantity int
}

// units returns the effective quantity, defaulting to 1 when unset.
func (i CartItem) units() int {
	if i.Quantity <= 0 {
		return 1
	}
	return i.Quantity
}

// subtotal returns the total cost for this line (amount times quantity).
func (i CartItem) subtotal() float64 {
	return i.Amount * float64(i.units())
}

// cart is the collection of items backing a Payment. It is unexported; callers
// interact with it through Payment.Add, Payment.Total and Payment.Info.
type cart struct {
	items []CartItem
}

// add appends an item to the cart.
func (c *cart) add(item CartItem) {
	c.items = append(c.items, item)
}

// total returns the sum of every line's subtotal.
func (c *cart) total() float64 {
	var total float64
	for _, item := range c.items {
		total += item.subtotal()
	}
	return total
}

// summary returns a comma-separated list of item titles, used as the
// "additionalinfo" field Paynow shows to the customer.
func (c *cart) summary() string {
	titles := make([]string, 0, len(c.items))
	for _, item := range c.items {
		titles = append(titles, item.Title)
	}
	return strings.Join(titles, ", ")
}
