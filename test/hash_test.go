package paynow_test

import (
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"

	"github.com/IamTyrone/paynow-go/internal/hash"
)

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		name           string
		values         url.Values
		integrationKey string
		wantNonEmpty   bool
	}{
		{
			name: "generates hash for valid payment data",
			values: func() url.Values {
				v := url.Values{}
				v.Set("id", "12345")
				v.Set("reference", "INV-001")
				v.Set("amount", "10.00")
				v.Set("authemail", "test@example.com")
				v.Set("phone", "0771234567")
				v.Set("method", "ecocash")
				v.Set("returnurl", "https://example.com/return")
				v.Set("resulturl", "https://example.com/result")
				v.Set("status", "Message")
				return v
			}(),
			integrationKey: "test-key-123",
			wantNonEmpty:   true,
		},
		{
			name: "generates consistent hash for same data",
			values: func() url.Values {
				v := url.Values{}
				v.Set("id", "12345")
				v.Set("amount", "50.00")
				return v
			}(),
			integrationKey: "secret-key",
			wantNonEmpty:   true,
		},
		{
			name:           "handles empty values",
			values:         url.Values{},
			integrationKey: "key",
			wantNonEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hash.GenerateHash(tt.values, tt.integrationKey)
			if tt.wantNonEmpty && got == "" {
				t.Error("GenerateHash() returned empty string, want non-empty")
			}
			if len(got) != 128 {
				t.Errorf("GenerateHash() returned hash of length %d, want 128 (SHA512 hex)", len(got))
			}
		})
	}
}

func TestGenerateHash_Consistency(t *testing.T) {
	values := url.Values{}
	values.Set("id", "12345")
	values.Set("reference", "INV-001")
	values.Set("amount", "10.00")
	integrationKey := "test-key"

	hash1 := hash.GenerateHash(values, integrationKey)
	hash2 := hash.GenerateHash(values, integrationKey)

	if hash1 != hash2 {
		t.Errorf("GenerateHash() not consistent: got %s and %s", hash1, hash2)
	}
}

func TestGenerateHash_ExcludesHashField(t *testing.T) {
	values1 := url.Values{}
	values1.Set("id", "12345")
	values1.Set("amount", "10.00")

	values2 := url.Values{}
	values2.Set("id", "12345")
	values2.Set("amount", "10.00")
	values2.Set("hash", "some-existing-hash")

	hash1 := hash.GenerateHash(values1, "key")
	hash2 := hash.GenerateHash(values2, "key")

	if hash1 != hash2 {
		t.Error("GenerateHash() should exclude 'hash' field from calculation")
	}
}

func TestGenerateHash_DifferentKeys(t *testing.T) {
	values := url.Values{}
	values.Set("id", "12345")
	values.Set("amount", "10.00")

	hash1 := hash.GenerateHash(values, "key1")
	hash2 := hash.GenerateHash(values, "key2")

	if hash1 == hash2 {
		t.Error("GenerateHash() should produce different hashes for different keys")
	}
}

func TestGenerateHash_DifferentData(t *testing.T) {
	values1 := url.Values{}
	values1.Set("amount", "10.00")

	values2 := url.Values{}
	values2.Set("amount", "20.00")

	hash1 := hash.GenerateHash(values1, "key")
	hash2 := hash.GenerateHash(values2, "key")

	if hash1 == hash2 {
		t.Error("GenerateHash() should produce different hashes for different data")
	}
}

func TestGenerateHash_UppercaseOutput(t *testing.T) {
	values := url.Values{}
	values.Set("id", "12345")

	h := hash.GenerateHash(values, "key")

	for _, c := range h {
		if c >= 'a' && c <= 'z' {
			t.Error("GenerateHash() should return uppercase hex string")
			break
		}
	}
}

func TestValidateHash(t *testing.T) {
	integrationKey := "test-integration-key"

	tests := []struct {
		name    string
		rawBody string
		wantErr bool
	}{
		{
			name:    "valid hash in response",
			rawBody: buildValidResponseBody(integrationKey),
			wantErr: false,
		},
		{
			name:    "invalid hash in response",
			rawBody: "status=Ok&browserurl=https://example.com&pollurl=https://poll.com&hash=INVALIDHASH",
			wantErr: true,
		},
		{
			name:    "missing hash in response",
			rawBody: "status=Ok&browserurl=https://example.com",
			wantErr: true,
		},
		{
			name:    "empty response",
			rawBody: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hash.ValidateHash(tt.rawBody, integrationKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHash_URLEncodedValues(t *testing.T) {
	integrationKey := "key123"

	var hashString string
	hashString += "Ok"
	hashString += "Payment successful"
	hashString += integrationKey

	importHash := hashGenerateSHA512(hashString)
	rawBody := "status=Ok&message=Payment+successful&hash=" + importHash

	err := hash.ValidateHash(rawBody, integrationKey)
	if err != nil {
		t.Errorf("ValidateHash() should handle URL-encoded values, got error: %v", err)
	}
}

func hashGenerateSHA512(input string) string {
	h := sha512.Sum512([]byte(input))
	return strings.ToUpper(hex.EncodeToString(h[:]))
}

func buildValidResponseBody(integrationKey string) string {
	status := "Ok"
	browserURL := "https://www.paynow.co.zw/payment/123"
	pollURL := "https://www.paynow.co.zw/interface/poll/123"

	hashInput := status + browserURL + pollURL + integrationKey
	hashBytes := sha512.Sum512([]byte(hashInput))
	h := strings.ToUpper(hex.EncodeToString(hashBytes[:]))

	return "status=" + status + "&browserurl=" + browserURL + "&pollurl=" + pollURL + "&hash=" + h
}
