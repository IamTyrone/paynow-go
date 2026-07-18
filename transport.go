package paynow

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// postForm sends body as an application/x-www-form-urlencoded POST to endpoint
// and returns the raw response body. The body is sent verbatim (not re-encoded)
// so field ordering — which is part of the Paynow hash — is preserved.
func (c *Client) postForm(ctx context.Context, endpoint, body string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("paynow: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("paynow: request to %s failed: %w", endpoint, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("paynow: failed to read response: %w", err)
	}
	return string(raw), nil
}
