# Paynow Zimbabwe Go SDK

A clean, idiomatic Go SDK for the [Paynow Zimbabwe](https://www.paynow.co.zw/) payment gateway, ported from the official [Node.js](https://github.com/paynow/Paynow-NodeJS-SDK) and [Python](https://github.com/paynow/Paynow-Python-SDK) SDKs.

- **Web payments** — redirect the customer to Paynow to choose how to pay.
- **Mobile (express checkout)** — charge EcoCash, OneMoney or InnBucks directly.
- **Status polling** and **result-URL webhooks**, with automatic hash verification.
- Context-aware, dependency-free (standard library only), and easy to mock in tests.

## Installation

```bash
go get github.com/IamTyrone/paynow-go
```

Requires Go 1.21 or newer.

## Getting started

Create a client with your integration credentials (available from your Paynow merchant dashboard):

```go
client := paynow.New(
    "YOUR_INTEGRATION_ID",
    "YOUR_INTEGRATION_KEY",
    paynow.WithResultURL("https://example.com/paynow/result"), // where Paynow POSTs status updates
    paynow.WithReturnURL("https://example.com/paynow/return"), // where the customer is sent back to
)
```

Build up a payment. A payment behaves like a shopping cart — add one or more items and the total is computed for you:

```go
payment := client.CreatePayment("INV-1001", "customer@example.com")
payment.Add("T-shirt", 10.00, 2) // title, unit amount, optional quantity (default 1)
payment.Add("Delivery", 3.50)
```

### Web payment

```go
resp, err := client.Send(context.Background(), payment)
if err != nil {
    log.Fatal(err)
}

// Redirect the customer to resp.RedirectURL to complete payment,
// and keep resp.PollURL to check the status later.
fmt.Println(resp.RedirectURL)
```

### Mobile (express checkout)

Mobile transactions charge the customer directly and require a valid auth email.

```go
resp, err := client.SendMobile(context.Background(), payment, "0771234567", paynow.MethodEcocash)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.PollURL)
if resp.Instructions != "" {
    fmt.Println(resp.Instructions) // e.g. a USSD prompt for the customer to confirm
}
```

Supported methods: `paynow.MethodEcocash`, `paynow.MethodOneMoney`, `paynow.MethodInnbucks`.

For InnBucks, the response carries the payment code, a deep link and a QR code:

```go
if resp.InnBucks != nil {
    fmt.Println("Code:", resp.InnBucks.AuthorizationCode)
    fmt.Println("Deep link:", resp.InnBucks.DeepLinkURL)
    fmt.Println("QR code:", resp.InnBucks.QRCodeURL)
}
```

## Checking transaction status

### Polling

```go
status, err := client.PollTransaction(context.Background(), resp.PollURL)
if err != nil {
    log.Fatal(err)
}

switch {
case status.Status.IsPaid():
    fmt.Println("Paid!")
case status.Status.IsFailed():
    fmt.Println("Failed.")
default:
    fmt.Println("Still pending:", status.Status)
}
```

### Result-URL webhook

When a transaction's status changes, Paynow POSTs a status update to your result URL. Pass the raw request body to `ProcessStatusUpdate`; the hash is verified for you:

```go
func resultHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)

    status, err := client.ProcessStatusUpdate(string(body))
    if err != nil {
        http.Error(w, "invalid update", http.StatusBadRequest)
        return
    }

    if status.Paid {
        // fulfil the order for status.Reference
    }
}
```

## Error handling

Network, parsing and hash-verification problems are returned as regular errors. Business errors reported by Paynow (for example an invalid integration id) are returned as `*paynow.APIError`, alongside a populated response you can still inspect:

```go
resp, err := client.Send(ctx, payment)
if err != nil {
    var apiErr *paynow.APIError
    if errors.As(err, &apiErr) {
        fmt.Println("Paynow said:", apiErr.Message)
    }
}
```

Sentinel errors you can match with `errors.Is`:

| Error | Meaning |
|-------|---------|
| `paynow.ErrNoPayment` | A nil payment was passed. |
| `paynow.ErrEmptyCart` | The payment has no items. |
| `paynow.ErrNonPositiveTotal` | The total is not greater than zero. |
| `paynow.ErrInvalidEmail` | A mobile payment lacks a valid auth email. |
| `paynow.ErrMissingHash` | A response that should be hashed had no hash. |
| `paynow.ErrHashMismatch` | A response hash did not match — possible tampering. |

## Custom HTTP client

By default the SDK uses a plain `*http.Client`. Supply your own (recommended, so you can set a timeout) with `WithHTTPClient`. Any type implementing `paynow.Doer` (which `*http.Client` satisfies) works, which also makes the SDK trivial to mock in tests:

```go
client := paynow.New(id, key,
    paynow.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

## Package layout

The public API lives in the root `paynow` package, split into small, focused files:

| File | Responsibility |
|------|----------------|
| `paynow.go` | `Client`, `New`, options |
| `payment.go`, `cart.go` | Building up a payment and its cart |
| `send.go` | `Send` / `SendMobile` and validation |
| `poll.go` | `PollTransaction` / `ProcessStatusUpdate` |
| `request.go`, `values.go` | Ordered request building and response parsing |
| `response.go` | `InitResponse` / `StatusResponse` / `InnBucksInfo` |
| `method.go`, `status.go` | Payment methods and transaction statuses |
| `errors.go` | Sentinel errors and `APIError` |
| `internal/hash` | SHA-512 request/response signing |

A complete, runnable flow lives in [`example/main.go`](example/main.go).

## License

See [LICENSE](LICENSE).
