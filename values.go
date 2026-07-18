package paynow

import (
	"net/url"
	"strings"

	"github.com/IamTyrone/paynow-go/internal/hash"
)

// orderedValues is an ordered set of key/value pairs. Order matters for Paynow's
// hashing scheme, so unlike url.Values (a map) this preserves insertion order
// both for building requests and for parsing responses.
type orderedValues struct {
	keys   []string
	values map[string]string
}

// newOrderedValues returns an empty orderedValues ready for use.
func newOrderedValues() *orderedValues {
	return &orderedValues{values: make(map[string]string)}
}

// set appends a key/value pair, preserving order. Keys are assumed unique, which
// holds for every Paynow request and response.
func (o *orderedValues) set(key, value string) {
	o.keys = append(o.keys, key)
	o.values[key] = value
}

// get returns the value for key and whether it was present.
func (o *orderedValues) get(key string) (string, bool) {
	v, ok := o.values[key]
	return v, ok
}

// signingValues returns every value except the hash field, in order, ready to be
// fed to the hashing routine.
func (o *orderedValues) signingValues() []string {
	out := make([]string, 0, len(o.keys))
	for _, k := range o.keys {
		if strings.EqualFold(k, "hash") {
			continue
		}
		out = append(out, o.values[k])
	}
	return out
}

// asMap returns the values as a plain map, exposed on responses via the Raw field.
func (o *orderedValues) asMap() map[string]string {
	m := make(map[string]string, len(o.values))
	for k, v := range o.values {
		m[k] = v
	}
	return m
}

// encode renders the values as an application/x-www-form-urlencoded body,
// preserving order. Standard library encoders sort keys, which would break the
// hash, so the body is assembled by hand here.
func (o *orderedValues) encode() string {
	var b strings.Builder
	for i, k := range o.keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(url.QueryEscape(k))
		b.WriteByte('=')
		b.WriteString(url.QueryEscape(o.values[k]))
	}
	return b.String()
}

// verifyHash checks the hash field against a hash computed from the other
// values. It returns ErrMissingHash when no hash is present and ErrHashMismatch
// when the hashes differ.
func (o *orderedValues) verifyHash(integrationKey string) error {
	received, ok := o.get("hash")
	if !ok {
		return ErrMissingHash
	}
	if !hash.Equal(received, o.signingValues(), integrationKey) {
		return ErrHashMismatch
	}
	return nil
}

// parseResponse parses a raw application/x-www-form-urlencoded body from Paynow
// into ordered, URL-decoded key/value pairs. Order is taken from the body so it
// can be used to reconstruct and verify the hash.
func parseResponse(body string) (*orderedValues, error) {
	body = strings.TrimSpace(body)
	body = strings.TrimPrefix(body, "?")

	ov := newOrderedValues()
	if body == "" {
		return ov, nil
	}

	for _, pair := range strings.Split(body, "&") {
		if pair == "" {
			continue
		}
		rawKey, rawValue, _ := strings.Cut(pair, "=")

		key, err := url.QueryUnescape(rawKey)
		if err != nil {
			return nil, err
		}
		value, err := url.QueryUnescape(rawValue)
		if err != nil {
			return nil, err
		}
		ov.set(key, value)
	}
	return ov, nil
}
