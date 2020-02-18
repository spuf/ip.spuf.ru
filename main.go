package main

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

type nope = struct{}

func main() {
	http.Handle("/", newHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func newHandler() http.Handler {
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

	dumper := func(r *http.Request) ([]byte, error) {
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
			return nil, err
		}
		var out bytes.Buffer
		out.WriteString(userIp)
		out.WriteString("\n\n")
		out.Write(bytes.ReplaceAll(dump, []byte("\r"), []byte("")))

		return out.Bytes(), nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		text, err := dumper(r)
		if err != nil {
			statusInternalServerError := http.StatusInternalServerError
			http.Error(w, http.StatusText(statusInternalServerError), statusInternalServerError)

			return
		}

		contentText := "text/plain"
		contentHtml := "text/html"
		contentType := selectContentType(r.Header.Get("Accept"), []string{contentHtml, contentText})
		if contentType == "" {
			contentType = contentText
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		var res bytes.Buffer
		if contentType == contentHtml {
			res.WriteString("<style>body{background-color:#fff;color:#000}@media screen and (prefers-color-scheme:dark){body{background-color:#000;color:#fff}}</style><pre style=word-wrap:break-word;white-space:pre-wrap>")
			template.HTMLEscape(&res, text)
			res.WriteString("</pre>")
		} else {
			res.Write(text)
		}

		w.Write(res.Bytes())
	})
}

func selectContentType(accept string, wanted []string) string {
	for _, part := range strings.Split(accept, ",") {
		for _, want := range wanted {
			if strings.HasPrefix(part, want) {
				return want
			}
		}
	}

	return ""
}
