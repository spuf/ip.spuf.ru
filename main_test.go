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
	if string(resBody) != `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body {
			background-color: white;
			color: black;
		}
		@media screen and (prefers-color-scheme: dark) {
			body{
				background-color: black;
				color: white;
			}
		}
		pre {
			word-wrap: break-word;
			white-space: pre-wrap;
		}
	</style>
</head>
<body>
	<pre>192.0.2.1:1234

GET / HTTP/1.1
Host: example.com
Accept: text/html

</pre>
</body>
</html>` {
		t.Fatalf("unexpected body: %s", resBody)
	}
}
