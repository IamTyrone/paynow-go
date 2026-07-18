package paynow_test

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// mockDoer is a paynow.Doer that captures the outgoing request and returns a
// canned response, for testing without hitting the network.
type mockDoer struct {
	response   string
	statusCode int
	err        error

	// capturedURL and capturedBody record the last request seen.
	capturedURL  string
	capturedBody string
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	m.capturedURL = req.URL.String()
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		m.capturedBody = string(body)
	}
	if m.err != nil {
		return nil, m.err
	}
	code := m.statusCode
	if code == 0 {
		code = http.StatusOK
	}
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(m.response)),
		Header:     make(http.Header),
	}, nil
}

// field is an ordered key/value pair used to build test response bodies.
type field struct{ key, value string }

// signResponse builds a URL-encoded response body from ordered fields and
// appends a valid Paynow hash computed with the same scheme the SDK uses.
func signResponse(key string, fields ...field) string {
	var hashInput strings.Builder
	var body strings.Builder
	for i, f := range fields {
		hashInput.WriteString(f.value)
		if i > 0 {
			body.WriteByte('&')
		}
		body.WriteString(url.QueryEscape(f.key))
		body.WriteByte('=')
		body.WriteString(url.QueryEscape(f.value))
	}
	hashInput.WriteString(strings.ToLower(key))

	sum := sha512.Sum512([]byte(hashInput.String()))
	h := strings.ToUpper(hex.EncodeToString(sum[:]))

	body.WriteString("&hash=")
	body.WriteString(h)
	return body.String()
}
