package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//go:embed templates/*.html
var templates embed.FS

func handleInput(w http.ResponseWriter, r *http.Request) string {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return ""
	}

	url := r.FormValue("long")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return ""
	}
	println(url)
	return url
}

func handleOutput(w http.ResponseWriter, r *http.Request, shortener *URLShortener) {
	longURL := handleInput(w, r)
	if longURL == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	shortCode := shortener.ShortUrl(longURL)
	fullShortURL := "http://localhost:8000/" + shortCode

	tmpl := template.Must(template.ParseFS(templates, "templates/index.html"))
	tmpl.Execute(w, map[string]string{"url": fullShortURL})
}

func redirect(w http.ResponseWriter, r *http.Request, shortener *URLShortener) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	originalURL := shortener.FindURL(shortCode)

	if originalURL == "" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}

func main() {
	shortener := &URLShortener{
		urls: make(map[string]string),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleOutput(w, r, shortener)
			return
		}

		if r.URL.Path == "/" { // ✅ Serve homepage for "/"
			tmpl := template.Must(template.ParseFS(templates, "templates/index.html"))
			tmpl.Execute(w, nil)
			return
		}

		redirect(w, r, shortener) // Handle short URLs
	})

	fmt.Println("Server started at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
