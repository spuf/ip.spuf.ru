package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/spuf/ip.spuf.ru/request_dumper"
)

func TestTextPlain(t *testing.T) {
	handler := newHandler(request_dumper.NewRequestDumper())

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "*/*")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatalf("unexpected status: %v", res.StatusCode)
	}
	contentType := res.Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Fatalf("unexpected content-type value: %v", contentType)
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	resBodyStr := string(resBody)
	if resBodyStr != "192.0.2.1:1234\n\nGET / HTTP/1.1\r\nHost: example.com\r\nAccept: */*\r\n\r\n\n" {
		t.Fatalf("unexpected body: %#v", resBodyStr)
	}
}

func TestTextHtml(t *testing.T) {
	handler := newHandler(request_dumper.NewRequestDumper())

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatalf("unexpected status: %v", res.StatusCode)
	}
	contentType := res.Header.Get("Content-Type")
	if contentType != "text/html" {
		t.Fatalf("unexpected content-type value: %v", contentType)
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	resBodyStr := string(resBody)
	if resBodyStr != "<!doctype html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <style>\n        body {\n            background-color: white;\n            color: black;\n        }\n        @media screen and (prefers-color-scheme: dark) {\n            body{\n                background-color: black;\n                color: white;\n            }\n        }\n        pre {\n            word-wrap: break-word;\n            white-space: pre-wrap;\n        }\n    </style>\n    <title>ip.spuf.ru</title>\n</head>\n<body>\n<pre>192.0.2.1:1234\n\nGET / HTTP/1.1\r\nHost: example.com\r\nAccept: text/html\r\n\r\n</pre>\n</body>\n</html>\n" {
		t.Fatalf("unexpected body: %#v", resBodyStr)
	}
}
