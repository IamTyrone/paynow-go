package paynow

import (
	"regexp"
	"strconv"
	"strings"
)

// emailPattern mirrors the validation used by the official Paynow SDKs.
var emailPattern = regexp.MustCompile(`^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$`)

// isValidEmail reports whether address looks like a valid email.
func isValidEmail(address string) bool {
	if address == "" {
		return false
	}
	return emailPattern.MatchString(address)
}

// formatAmount renders a monetary amount with two decimal places, the format
// Paynow expects.
func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}

// equalFoldTrim reports whether a and b are equal ignoring case and surrounding
// whitespace. Paynow is inconsistent with casing on status fields.
func equalFoldTrim(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), b)
}
