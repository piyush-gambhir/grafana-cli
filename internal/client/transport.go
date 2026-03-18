package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// verboseTransport wraps an http.RoundTripper and logs requests/responses to stderr.
type verboseTransport struct {
	inner  http.RoundTripper
	output io.Writer
}

func (t *verboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log request.
	fmt.Fprintf(t.output, "--> %s %s\n", req.Method, req.URL.String())

	// Log headers, redacting Authorization.
	for key, values := range req.Header {
		for _, v := range values {
			if key == "Authorization" {
				fmt.Fprintf(t.output, "    %s: [REDACTED]\n", key)
			} else {
				fmt.Fprintf(t.output, "    %s: %s\n", key, v)
			}
		}
	}

	resp, err := t.inner.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(t.output, "<-- ERROR: %v (%s)\n", err, time.Since(start).Round(time.Millisecond))
		return nil, err
	}

	fmt.Fprintf(t.output, "<-- %d %s (%s)\n", resp.StatusCode, http.StatusText(resp.StatusCode), time.Since(start).Round(time.Millisecond))
	return resp, nil
}

// EnableVerboseLogging wraps the client's HTTP transport with verbose logging.
func (c *Client) EnableVerboseLogging(w io.Writer) {
	transport := c.HTTPClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	c.HTTPClient.Transport = &verboseTransport{inner: transport, output: w}
}
