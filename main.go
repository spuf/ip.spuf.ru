package main

import (
	"fmt"
	"net/http"
	"sort"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", handler)
	appengine.Main()
}

func handler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprint(res, RequestAsString(req))
}

var hiddenHeaders = map[string]bool{
	"X-Appengine-Default-Namespace": true,
	"X-Cloud-Trace-Context":         true,
	"X-Google-Apps-Metadata":        true,
	"X-Zoo":                         true,
}

// RequestAsString generates text representation of HTTP request
func RequestAsString(req *http.Request) string {
	var res string
	res += fmt.Sprintf("%s\n\n", req.RemoteAddr)
	res += fmt.Sprintf("%s %s %s\n", req.Method, req.RequestURI, req.Proto)

	var names []string
	for name := range req.Header {
		if _, ok := hiddenHeaders[name]; !ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		for _, value := range req.Header[name] {
			res += fmt.Sprintf("%s: %s\n", name, value)
		}
	}
	return res
}
