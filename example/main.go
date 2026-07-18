// Command example demonstrates a full Paynow payment flow with the Go SDK:
// creating a payment, initiating it via express-checkout mobile money, and
// polling for the final status.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IamTyrone/paynow-go"
)

func main() {
	client := paynow.New(
		"YOUR_INTEGRATION_ID",
		"YOUR_INTEGRATION_KEY",
		paynow.WithResultURL("https://example.com/paynow/result"),
		paynow.WithReturnURL("https://example.com/paynow/return"),
		paynow.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	)

	// Build up the payment (acts like a shopping cart).
	payment := client.CreatePayment("INV-1001", "customer@example.com")
	payment.Add("T-shirt", 10.00, 2)
	payment.Add("Delivery", 3.50)

	ctx := context.Background()

	// Initiate an EcoCash express-checkout payment.
	resp, err := client.SendMobile(ctx, payment, "0771234567", paynow.MethodEcocash)
	if err != nil {
		var apiErr *paynow.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("Paynow rejected the payment: %s", apiErr.Message)
		}
		log.Fatalf("Failed to initiate payment: %v", err)
	}

	fmt.Println("Payment initiated.")
	fmt.Println("Poll URL:", resp.PollURL)
	if resp.Instructions != "" {
		fmt.Println("Instructions:", resp.Instructions)
	}

	// Poll until the transaction reaches a terminal state.
	for {
		status, err := client.PollTransaction(ctx, resp.PollURL)
		if err != nil {
			log.Fatalf("Failed to poll transaction: %v", err)
		}

		fmt.Println("Status:", status.Status)
		switch {
		case status.Status.IsPaid():
			fmt.Println("Payment successful!")
			return
		case status.Status.IsFailed():
			fmt.Println("Payment failed.")
			return
		}

		time.Sleep(5 * time.Second)
	}
}
