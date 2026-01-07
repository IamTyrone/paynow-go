# Paynow Zimbabwe Go SDK

A Go SDK for integrating with the [Paynow Zimbabwe](https://www.paynow.co.zw/) payment gateway.

## Installation

```bash
go get github.com/IamTyrone/paynow-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

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

    // Send a mobile payment (EcoCash)
    response, err := client.SendMobile(paynow.Payment{
        Reference: "INV-1001",
        Amount:    10.00,
        AuthEmail: "customer@example.com",
        Phone:     "0771234567",
        Method:    types.MethodEcocash,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %s\n", response.Status)
    fmt.Printf("Poll URL: %s\n", response.PollURL)
}
```

## Configuration

| Field            | Description                                         |
| ---------------- | --------------------------------------------------- |
| `IntegrationID`  | Your Paynow integration ID                          |
| `IntegrationKey` | Your Paynow integration key                         |
| `ResultURL`      | URL where Paynow will post transaction results      |
| `ReturnURL`      | URL where the user will be redirected after payment |

You can find your integration details in your [Paynow Dashboard](https://www.paynow.co.zw/).

## Payment Methods

Currently supported payment methods:

| Method  | Constant              |
| ------- | --------------------- |
| EcoCash | `types.MethodEcocash` |

## API Reference

### Creating a Client

```go
client := paynow.New(paynow.Config{
    IntegrationID:  "your-id",
    IntegrationKey: "your-key",
    ResultURL:      "https://example.com/result",
    ReturnURL:      "https://example.com/return",
})
```

### Sending a Mobile Payment

```go
response, err := client.SendMobile(paynow.Payment{
    Reference: "INV-1001",
    Amount:    10.00,
    AuthEmail: "customer@example.com",
    Phone:     "0771234567",
    Method:    types.MethodEcocash,
})
```

**Payment Fields:**

| Field       | Type            | Required | Description                          |
| ----------- | --------------- | -------- | ------------------------------------ |
| `Reference` | `string`        | Yes      | Unique transaction reference         |
| `Amount`    | `float64`       | Yes      | Payment amount                       |
| `AuthEmail` | `string`        | Yes      | Customer's email address             |
| `Phone`     | `string`        | Yes      | Customer's phone number              |
| `Method`    | `PaymentMethod` | No       | Payment method (defaults to EcoCash) |

**Response Fields:**

| Field        | Type     | Description                        |
| ------------ | -------- | ---------------------------------- |
| `Status`     | `string` | Transaction status                 |
| `BrowserURL` | `string` | URL for browser-based payment      |
| `PollURL`    | `string` | URL to poll for transaction status |
| `Hash`       | `string` | Response hash                      |
| `Error`      | `string` | Error message (if any)             |

### Polling Transaction Status

```go
status, err := client.PollTransaction(response.PollURL)
if err != nil {
    log.Fatal(err)
}

if status.Status.IsPaid() {
    fmt.Println("Payment successful!")
}
```

**Status Response Fields:**

| Field             | Type                | Description                |
| ----------------- | ------------------- | -------------------------- |
| `Reference`       | `string`            | Your transaction reference |
| `Amount`          | `float64`           | Transaction amount         |
| `PaynowReference` | `string`            | Paynow's reference number  |
| `PollURL`         | `string`            | URL for future polling     |
| `Status`          | `TransactionStatus` | Current transaction status |
| `Hash`            | `string`            | Response hash              |

### Transaction Status Helpers

```go
status.Status.IsPaid()    // Returns true if payment was successful
status.Status.IsPending() // Returns true if payment is still pending
status.Status.IsFailed()  // Returns true if payment failed or was cancelled
```

**Possible Status Values:**

| Status      | Description                  |
| ----------- | ---------------------------- |
| `Created`   | Transaction created          |
| `Sent`      | Transaction sent to provider |
| `Pending`   | Awaiting payment             |
| `Paid`      | Payment successful           |
| `Cancelled` | Transaction cancelled        |
| `Failed`    | Transaction failed           |
| `Refunded`  | Transaction refunded         |

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/IamTyrone/paynow-go"
    "github.com/IamTyrone/paynow-go/types"
)

func main() {
    client := paynow.New(paynow.Config{
        IntegrationID:  "YOUR_INTEGRATION_ID",
        IntegrationKey: "YOUR_INTEGRATION_KEY",
        ResultURL:      "https://example.com/result",
        ReturnURL:      "https://example.com/return",
    })

    // Initiate payment
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

    fmt.Printf("Payment initiated! Poll URL: %s\n", response.PollURL)

    // Poll for status
    for {
        status, err := client.PollTransaction(response.PollURL)
        if err != nil {
            log.Fatalf("Failed to poll: %v", err)
        }

        fmt.Printf("Status: %s\n", status.Status)

        if status.Status.IsPaid() {
            fmt.Printf("Payment successful! Paynow Ref: %s\n", status.PaynowReference)
            break
        }

        if status.Status.IsFailed() {
            fmt.Println("Payment failed or cancelled")
            break
        }

        time.Sleep(5 * time.Second)
    }
}
```

## Testing

All tests are located in the dedicated `test/` folder for better organization.

Run all tests:

```bash
go test -v ./test
```

Run tests with coverage:

```bash
go test -cover ./test
```

Run tests with verbose output:

```bash
go test -v ./test
```

### Test Coverage

The SDK maintains high test coverage:

| Package         | Coverage |
| --------------- | -------- |
| `paynow`        | ~94%     |
| `internal/hash` | ~97%     |
| `types`         | 100%     |

### Mocking for Tests

The SDK provides an `HTTPClient` interface for mocking HTTP requests in your tests:

```go
type MockHTTPClient struct {
    PostFormFunc func(url string, data url.Values) (*http.Response, error)
    GetFunc      func(url string) (*http.Response, error)
}

func (m *MockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
    return m.PostFormFunc(url, data)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
    return m.GetFunc(url)
}

// Usage
mockClient := &MockHTTPClient{
    PostFormFunc: func(url string, data url.Values) (*http.Response, error) {
        // Return mock response
        return &http.Response{
            StatusCode: 200,
            Body:       io.NopCloser(bytes.NewBufferString("status=Ok&hash=...")),
        }, nil
    },
}

client := paynow.NewWithHTTPClient(config, mockClient)
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Note:** All PRs must pass the automated test suite before merging. Tests run automatically on every PR via GitHub Actions.

## License

MIT License - see [LICENSE](LICENSE) for details.
