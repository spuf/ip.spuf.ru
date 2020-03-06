package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/spuf/ip.spuf.ru/request_dumper"
)

func main() {
	http.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(w, r, "./assets/favicon.ico")
	}))

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
	templates := template.Must(template.ParseFiles("./assets/index.html", "./assets/index.txt"))
	contentText := "text/plain"
	contentHtml := "text/html"
	templatesMap := map[string]string{
		contentText: "index.txt",
		contentHtml: "index.html",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		requestDump, err := dumper(r)
		if err != nil {
			panic(err)
		}

		contentType := selectContentType(r.Header.Get("Accept"), []string{contentHtml, contentText}, contentText)

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		if err := templates.ExecuteTemplate(w, templatesMap[contentType], requestDump); err != nil {
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
