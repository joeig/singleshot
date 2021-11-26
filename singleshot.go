// Package singleshot provides an http.RoundTripper which deduplicates similar HTTP requests.
//
// If two similar HTTP requests are supposed to be sent concurrently, the first one will actually be sent to the server, while the second one waits until the first one was fulfilled completely.
// The second request will never be sent to the server, but returns a copy of the response of the first request.
//
//  Req 1 -----------------> Resp 1
//               Req 2 ----> Resp 1'
//                                  Req 3 -------------> Resp 3
package singleshot

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
	"strings"

	"golang.org/x/sync/singleflight"
)

// Transport is safe for concurrent use.
type Transport struct {
	transport    http.RoundTripper
	requestGroup singleflight.Group
}

// NewTransport creates a new instance of singleshot.Transport.
func NewTransport(transport http.RoundTripper) *Transport {
	return &Transport{
		transport:    transport,
		requestGroup: singleflight.Group{},
	}
}

// RoundTrip deduplicates similar subsequential HTTP requests, if the first request of a kind
// has not been completely fulfilled yet.
//
// Only "GET" requests (excluding "range" requests) are deduplicated, other requests are passed.
// Request are considered similar, if they share the same method and request URI (see RFC 2616, section 9.3).
//
// Always apply proper timeouts or use requests with contexts, otherwise one request which is timing out
// may stop subsequent requests from being retried.
//
// While the response body is a valid io.ReadCloser, the transfer itself has finished.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !isDeduplicatable(req) {
		return t.transport.RoundTrip(req)
	}

	respBytes, err, _ := t.requestGroup.Do(groupKey(req), func() (interface{}, error) {
		resp, err := t.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		return httputil.DumpResponse(resp, true)
	})

	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(respBytes.([]byte))
	return http.ReadResponse(bufio.NewReader(b), req)
}

func isDeduplicatable(req *http.Request) bool {
	const rangeHeader = "range"
	return req.Method == http.MethodGet && req.Header.Get(rangeHeader) == ""
}

func groupKey(req *http.Request) string {
	return strings.Join([]string{req.Method, req.URL.String()}, " ")
}
