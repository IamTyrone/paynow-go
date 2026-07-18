package paynow_test

import (
	"testing"

	"github.com/IamTyrone/paynow-go"
)

func TestPayment_TotalAndInfo(t *testing.T) {
	p := paynow.NewPayment("INV-1", "buyer@example.com").
		Add("Widget", 5.50, 2).
		Add("Sticker", 1.00)

	if got := p.Total(); got != 12.00 {
		t.Errorf("Total() = %.2f, want 12.00", got)
	}
	if got := p.Info(); got != "Widget, Sticker" {
		t.Errorf("Info() = %q, want %q", got, "Widget, Sticker")
	}
	if got := len(p.Items()); got != 2 {
		t.Errorf("len(Items()) = %d, want 2", got)
	}
}

func TestPayment_DefaultQuantity(t *testing.T) {
	p := paynow.NewPayment("INV-2", "buyer@example.com").Add("Book", 20.00)
	if got := p.Total(); got != 20.00 {
		t.Errorf("Total() = %.2f, want 20.00 (quantity should default to 1)", got)
	}
}

func TestPayment_ItemsIsCopy(t *testing.T) {
	p := paynow.NewPayment("INV-3", "buyer@example.com").Add("Pen", 2.00)
	items := p.Items()
	items[0].Amount = 999

	if p.Total() != 2.00 {
		t.Error("mutating the slice from Items() should not affect the payment")
	}
}

func TestClient_CreatePayment(t *testing.T) {
	client := paynow.New("id", "key")
	p := client.CreatePayment("INV-4", "buyer@example.com")

	if p.Reference != "INV-4" || p.AuthEmail != "buyer@example.com" {
		t.Errorf("CreatePayment() = %+v, want reference/authEmail set", p)
	}
}
