// Package hash provides hash generation and validation utilities for Paynow API requests.
package hash

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// GenerateHash creates a SHA512 hash for the request data.
// Keys are sorted alphabetically to match http.PostForm behavior.
func GenerateHash(values url.Values, integrationKey string) string {
	var keys []string
	for k := range values {
		if strings.ToLower(k) != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var hashString strings.Builder
	for _, k := range keys {
		hashString.WriteString(values.Get(k))
	}
	hashString.WriteString(integrationKey)

	hash := sha512.Sum512([]byte(hashString.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// ValidateHash validates the hash in a response from Paynow.
// It parses the raw body to preserve the order of fields.
func ValidateHash(rawBody string, integrationKey string) error {
	pairs := strings.Split(rawBody, "&")
	var hashString strings.Builder
	var receivedHash string

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]

		if strings.ToLower(key) == "hash" {
			receivedHash = value
			continue
		}

		decodedValue, err := url.QueryUnescape(value)
		if err != nil {
			return fmt.Errorf("failed to decode value for key %s: %w", key, err)
		}
		hashString.WriteString(decodedValue)
	}

	hashString.WriteString(integrationKey)
	hash := sha512.Sum512([]byte(hashString.String()))
	expectedHash := strings.ToUpper(hex.EncodeToString(hash[:]))

	if receivedHash != expectedHash {
		return fmt.Errorf("invalid hash: received %s, expected %s", receivedHash, expectedHash)
	}
	return nil
}
