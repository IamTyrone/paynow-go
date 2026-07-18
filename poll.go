package paynow

import "context"

// PollTransaction checks the current status of a transaction using the poll URL
// returned when the transaction was initiated. The response hash is verified for
// non-error responses.
func (c *Client) PollTransaction(ctx context.Context, pollURL string) (*StatusResponse, error) {
	raw, err := c.postForm(ctx, pollURL, "")
	if err != nil {
		return nil, err
	}

	values, err := parseResponse(raw)
	if err != nil {
		return nil, err
	}

	status, _ := values.get("status")
	if equalFoldTrim(status, responseError) {
		resp := newStatusResponse(values)
		return resp, &APIError{Message: resp.Error}
	}

	if err := values.verifyHash(c.integrationKey); err != nil {
		return nil, err
	}
	return newStatusResponse(values), nil
}

// ProcessStatusUpdate parses and verifies a status update that Paynow posts to
// your result URL. Pass the raw request body (for example the bytes read from
// http.Request.Body) so the hash can be verified against the exact field order
// Paynow used.
func (c *Client) ProcessStatusUpdate(rawBody string) (*StatusResponse, error) {
	values, err := parseResponse(rawBody)
	if err != nil {
		return nil, err
	}

	status, _ := values.get("status")
	if equalFoldTrim(status, responseError) {
		resp := newStatusResponse(values)
		return resp, &APIError{Message: resp.Error}
	}

	if err := values.verifyHash(c.integrationKey); err != nil {
		return nil, err
	}
	return newStatusResponse(values), nil
}
