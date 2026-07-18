package paynow

// API endpoints used by the SDK to talk to Paynow.
const (
	// urlInitiateTransaction is the endpoint for normal, web-based
	// transactions where the customer is redirected to Paynow to pay.
	urlInitiateTransaction = "https://www.paynow.co.zw/interface/initiatetransaction"

	// urlInitiateMobileTransaction is the endpoint for express-checkout
	// mobile money transactions (EcoCash, OneMoney, InnBucks, ...).
	urlInitiateMobileTransaction = "https://www.paynow.co.zw/interface/remotetransaction"
)

// Status values returned by Paynow in the "status" field of a response.
const (
	// responseOK indicates the request was accepted by Paynow.
	responseOK = "ok"

	// responseError indicates Paynow rejected the request. Error responses
	// are not hashed, so hash verification is skipped for them.
	responseError = "error"
)

// InnBucks helper prefixes. When Paynow returns an authorization code for an
// InnBucks payment, these are combined with that code to build a deep link the
// customer can tap and a QR code they can scan to complete the payment.
const (
	innbucksDeepLinkPrefix = "schinn.wbpycode://innbucks.co.zw?pymInnCode="
	googleQRPrefix         = "https://chart.googleapis.com/chart?cht=qr&chs=200x200&chl="
)
