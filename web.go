package web

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprintf(res, "%s\n\n", req.RemoteAddr)
	fmt.Fprintf(res, "%s %s %s\n", req.Method, req.RequestURI, req.Proto)
	fmt.Fprintf(res, "Host: %s\n", req.Host)
	var names []string
	for name := range req.Header {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		for _, value := range req.Header[name] {
			fmt.Fprintf(res, "%s: %s\n", name, value)
		}
	}
}
