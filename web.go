package web

import (
	"fmt"
	"net/http"
	"sort"
)

var hiddenHeaders = map[string]bool{
	"X-Cloud-Trace-Context":  true,
	"X-Google-Apps-Metadata": true,
	"X-Zoo":                  true,
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprintf(res, "%s\n\n", req.RemoteAddr)
	fmt.Fprintf(res, "%s %v %s\n", req.Method, req.URL, req.Proto)
	var names []string
	for name := range req.Header {
		if _, ok := hiddenHeaders[name]; !ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		for _, value := range req.Header[name] {
			fmt.Fprintf(res, "%s: %s\n", name, value)
		}
	}
}
