# Singleshot

Singleshot provides an `http.RoundTripper` which deduplicates similar HTTP requests.

[![Go Report Card](https://goreportcard.com/badge/go.eigsys.de/singleshot)](https://goreportcard.com/report/go.eigsys.de/singleshot)
[![PkgGoDev](https://pkg.go.dev/badge/go.eigsys.de/singleshot)](https://pkg.go.dev/go.eigsys.de/singleshot)

If two similar HTTP requests are supposed to be sent concurrently, the first one will actually be sent to the server, while the second one waits until the first one was fulfilled completely.
The second request will never be sent to the server, but returns a copy of the response of the first request.

```mermaid
sequenceDiagram
    participant Client
    participant Server
    
    Client ->> +Server: Request 1
    Note over Client: Identical, subsequent requests wait<br/>until Request 1 was fulfilled.
    Client -> Client: Request 2
    Server -->> -Client: Response 1, Response 2
```

## Usage

```go
package main

import (
	"net/http"

	"go.eigsys.de/singleshot"
)

func main() {
	_ = http.Client{
		Transport: singleshot.NewTransport(http.DefaultTransport),
	}
}
```

## Notes

* Always apply proper timeouts or use requests with contexts, otherwise one request which is timing out may stop subsequent requests from being retried.
* Requests are considered deduplicatable, if they share the same HTTP method and request URI. Furthermore, the method has to be `GET` and the request must not be a `range` request. The body is ignored according to [RFC 2616, section 9.3](https://www.rfc-editor.org/rfc/rfc2616#section-9.3).

## Documentation

See [pkg.go.dev](https://pkg.go.dev/go.eigsys.de/singleshot).
