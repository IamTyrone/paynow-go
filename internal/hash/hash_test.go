package hash_test

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/IamTyrone/paynow-go/internal/hash"
)

func reference(values []string, key string) string {
	sum := sha512.Sum512([]byte(strings.Join(values, "") + strings.ToLower(key)))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func TestMake_MatchesScheme(t *testing.T) {
	values := []string{"10.00", "INV-1", "12345"}
	got := hash.Make(values, "SECRET-Key")
	want := reference(values, "SECRET-Key")

	if got != want {
		t.Fatalf("Make() = %s, want %s", got, want)
	}
	if len(got) != 128 {
		t.Errorf("Make() length = %d, want 128", len(got))
	}
}

func TestMake_LowercasesKey(t *testing.T) {
	values := []string{"a", "b"}
	if hash.Make(values, "KEY") != hash.Make(values, "key") {
		t.Error("Make() should lower-case the integration key before hashing")
	}
}

func TestMake_IsUppercaseHex(t *testing.T) {
	got := hash.Make([]string{"x"}, "key")
	if got != strings.ToUpper(got) {
		t.Errorf("Make() = %s, want uppercase", got)
	}
}

func TestMake_OrderMatters(t *testing.T) {
	if hash.Make([]string{"a", "b"}, "k") == hash.Make([]string{"b", "a"}, "k") {
		t.Error("Make() should depend on value order")
	}
}

func TestEqual(t *testing.T) {
	values := []string{"1", "2"}
	h := hash.Make(values, "key")

	if !hash.Equal(h, values, "key") {
		t.Error("Equal() should match a freshly generated hash")
	}
	if !hash.Equal(strings.ToLower(h), values, "key") {
		t.Error("Equal() should be case-insensitive")
	}
	if hash.Equal("nope", values, "key") {
		t.Error("Equal() should reject a wrong hash")
	}
}
