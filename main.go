package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"

	"github.com/spuf/ip.spuf.ru/request_dumper"
	"github.com/spuf/ip.spuf.ru/websocket_ping"
)

func main() {
	dumper := request_dumper.NewRequestDumper()
	http.Handle("/", newHandler(dumper))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func newHandler(dumper request_dumper.RequestDumper) http.Handler {
	defaultContentType := "text/plain"
	var templates *template.Template
	if err := pkger.Walk("/static/templates", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			return fmt.Errorf("unknow content type for %v", path)
		}
		p := strings.SplitN(contentType, ";", 2)
		contentType = p[0]
		f, err := pkger.Open(path)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		var tmpl *template.Template
		if templates == nil {
			templates = template.New(contentType)
		}
		if contentType == templates.Name() {
			tmpl = templates
		} else {
			tmpl = templates.New(contentType)
		}
		if _, err := tmpl.Parse(string(content)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}

	if t := templates.Lookup(defaultContentType); t == nil {
		panic(fmt.Errorf("default content type %v not found", defaultContentType))
	}
	var availableContentTypes []string
	for _, t := range templates.Templates() {
		availableContentTypes = append(availableContentTypes, t.Name())
	}

	fs := http.FileServer(pkger.Dir("/static/public"))
	ws := websocket_ping.NewWebsocketPing()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			ws.ServeHTTP(w, r)
			return
		}

		if r.URL.Path != "/" {
			fs.ServeHTTP(w, r)
			return
		}

		requestDump, err := dumper(r)
		if err != nil {
			panic(err)
		}

		contentType := selectContentType(r.Header.Get("Accept"), availableContentTypes, defaultContentType)

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		if err := templates.ExecuteTemplate(w, contentType, requestDump); err != nil {
			panic(err)
		}
	})
}

func selectContentType(accept string, wanted []string, fallback string) string {
	for _, part := range strings.Split(accept, ",") {
		for _, want := range wanted {
			if strings.HasPrefix(part, want) {
				return want
			}
		}
	}

	return fallback
}
