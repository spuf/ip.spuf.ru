package request_dumper

import (
	"net/http"
	"net/http/httputil"
	"strings"
)

type nope struct{}

type RequestDump struct {
	UserIp  string
	Request []byte
}

type RequestDumper func(*http.Request) (*RequestDump, error)

func NewRequestDumper() RequestDumper {
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

	return func(r *http.Request) (*RequestDump, error) {
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

		requestDump := &RequestDump{
			UserIp:  userIp,
			Request: dump,
		}

		return requestDump, nil
	}
}
