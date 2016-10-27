package web

import (
	"net/http"
	"testing"
)

func TestRequestAsString(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "127.0.0.1",
		Method:     "GET",
		RequestURI: "/test",
		Proto:      "HTTP/1.1",
		Header: http.Header{
			"X-Zoo": {"Empty"},
			"Dnt":   {"1"},
		},
	}
	actual := RequestAsString(req)
	expected := `127.0.0.1

GET /test HTTP/1.1
Dnt: 1
`
	if actual != expected {
		t.Errorf("%#v != %#v", actual, expected)
	}
}
