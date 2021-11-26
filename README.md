# Singleshot

Singleshot provides an `http.RoundTripper` which deduplicates similar HTTP requests.

[![Build Status](https://github.com/joeig/singleshot/workflows/Tests/badge.svg)](https://github.com/joeig/singleshot/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeig/singleshot)](https://goreportcard.com/report/github.com/joeig/singleshot)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/joeig/singleshot)](https://pkg.go.dev/github.com/joeig/singleshot)

If two similar HTTP requests are supposed to be sent concurrently, the first one will actually be sent to the server, while the second one waits until the first one was fulfilled completely.
The second request will never be sent to the server, but returns a copy of the response of the first request.

```text
Req 1 -----------------> Resp 1
             Req 2 ----> Resp 1'
                                Req 3 -------------> Resp 3
```

## Usage

```go
import (
	"github.com/joeig/singleshot"
)

_ = http.Client{
	Transport: singleshot.NewTransport(http.DefaultTransport),
}
```

## Notes

* Always apply proper timeouts or use requests with contexts, otherwise one request which is timing out may stop subsequent requests from being retried.
* Requests are considered deduplicatable, if they share the same HTTP method and request URI. Furthermore, the method has to be `GET` and the request must not be a `range` request. The body is ignored according to [RFC 2616, section 9.3](https://www.rfc-editor.org/rfc/rfc2616#section-9.3).

## Documentation

See [GoDoc](https://godoc.org/github.com/joeig/singleshot).
