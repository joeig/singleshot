package singleshot_test

import (
	"net/http"
	"testing"

	"github.com/joeig/singleshot"
)

func TestNewTransport(t *testing.T) {
	var _ http.RoundTripper = (*singleshot.Transport)(nil)
}

func ExampleNewTransport() {
	_ = http.Client{
		Transport: singleshot.NewTransport(http.DefaultTransport),
	}
}
