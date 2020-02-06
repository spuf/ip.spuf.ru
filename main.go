package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

type nope = struct{}

func main() {
	hiddenHeaders := map[string]nope{
		"Accept-Encoding":        nope{},
		"Connection":             nope{},
		"Forwarded":              nope{},
		"Keep-Alive":             nope{},
		"Proxy-Authorization":    nope{},
		"TE":                     nope{},
		"Trailer":                nope{},
		"Transfer-Encoding":      nope{},
		"X-Cloud-Trace-Context":  nope{},
		"X-Forwarded-For":        nope{},
		"X-Forwarded-Proto":      nope{},
		"X-Google-Apps-Metadata": nope{},
	}
	appengineHeaderPrefix := "X-Appengine-"
	allowedAppengineHeaders := map[string]nope{
		appengineHeaderPrefix + "City":    nope{},
		appengineHeaderPrefix + "Country": nope{},
		appengineHeaderPrefix + "Region":  nope{},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		safeHeader := http.Header{}
		userIp := r.RemoteAddr
		for key, values := range r.Header {
			if key == "X-Appengine-User-Ip" {
				userIp = values[0]

				continue
			}
			if _, ok := hiddenHeaders[key]; ok {
				continue
			}
			if strings.HasPrefix(key, appengineHeaderPrefix) {
				if _, ok := allowedAppengineHeaders[key]; !ok {
					continue
				}
			}

			for _, value := range values {
				safeHeader.Add(key, value)
			}
		}
		r.Header = safeHeader

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			statusInternalServerError := http.StatusInternalServerError
			http.Error(w, http.StatusText(statusInternalServerError), statusInternalServerError)

			return
		}

		io.WriteString(w, userIp)
		io.WriteString(w, "\n\n")
		w.Write(dump)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
