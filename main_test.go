package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestTextPlain(t *testing.T) {
	handler := newHandler()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "*/*")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", res.StatusCode)
	}
	contentType := res.Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Fatalf("unexpected content-type value: %s", contentType)
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(resBody) != "192.0.2.1:1234\n\nGET / HTTP/1.1\nHost: example.com\nAccept: */*\n\n" {
		t.Fatalf("unexpected body: %s", resBody)
	}
}

func TestTextHtml(t *testing.T) {
	handler := newHandler()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", res.StatusCode)
	}
	contentType := res.Header.Get("Content-Type")
	if contentType != "text/html" {
		t.Fatalf("unexpected content-type value: %s", contentType)
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(resBody) != "<style>body{background-color:#fff;color:#000}@media screen and (prefers-color-scheme:dark){body{background-color:#000;color:#fff}}</style><pre style=word-wrap:break-word;white-space:pre-wrap>192.0.2.1:1234\n\nGET / HTTP/1.1\nHost: example.com\nAccept: text/html\n\n</pre>" {
		t.Fatalf("unexpected body: %s", resBody)
	}
}
