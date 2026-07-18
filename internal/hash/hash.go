// Package hash implements the SHA-512 message hashing scheme Paynow uses to
// sign and verify every request and response.
//
// The scheme is: concatenate all field values (in order, excluding any existing
// "hash" field), append the integration key, take the SHA-512 digest and encode
// it as an uppercase hex string. Paynow lower-cases the integration key before
// hashing, so this package does too.
package hash

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

// Make computes the Paynow hash for the given ordered values and integration
// key. Values must be supplied in the same order they are transmitted to (or
// were received from) Paynow, since ordering is part of the signature.
func Make(values []string, integrationKey string) string {
	var b strings.Builder
	for _, v := range values {
		b.WriteString(v)
	}
	b.WriteString(strings.ToLower(integrationKey))

	sum := sha512.Sum512([]byte(b.String()))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// Equal reports whether the hash computed from values and integrationKey
// matches the provided candidate hash. The comparison is case-insensitive.
func Equal(candidate string, values []string, integrationKey string) bool {
	return strings.EqualFold(candidate, Make(values, integrationKey))
}
