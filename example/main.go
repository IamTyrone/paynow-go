// Example demonstrates how to use the Paynow Go SDK.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IamTyrone/paynow-go"
	"github.com/IamTyrone/paynow-go/types"
)

func main() {
	// Create a new Paynow client
	client := paynow.New(paynow.Config{
		IntegrationID:  "YOUR_INTEGRATION_ID",
		IntegrationKey: "YOUR_INTEGRATION_KEY",
		ResultURL:      "https://example.com/paynow/result",
		ReturnURL:      "https://example.com/paynow/return",
	})

	// Create and send a mobile payment
	response, err := client.SendMobile(paynow.Payment{
		Reference: "INV-1001",
		Amount:    10.00,
		AuthEmail: "customer@example.com",
		Phone:     "0771234567",
		Method:    types.MethodEcocash,
	})
	if err != nil {
		log.Fatalf("Failed to initiate payment: %v", err)
	}

	fmt.Printf("Payment initiated successfully!\n")
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Poll URL: %s\n", response.PollURL)

	// Poll for transaction status
	for {
		status, err := client.PollTransaction(response.PollURL)
		if err != nil {
			log.Fatalf("Failed to poll transaction: %v", err)
		}

		fmt.Printf("Transaction status: %s\n", status.Status)

		if status.Status.IsPaid() {
			fmt.Println("Payment successful!")
			break
		}

		if status.Status.IsFailed() {
			fmt.Println("Payment failed!")
			break
		}

		// Wait before polling again
		time.Sleep(5 * time.Second)
	}
}
