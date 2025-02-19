package controllers

import (
	"database/sql"
	"github.com/irongollem/urlShortner.git/internal/db"
	"html/template"
	"math/big"
	"net/http"
	"strings"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func shortenUrl(url string) string {
	n := new(big.Int).SetBytes([]byte(strings.ToLower(url)))
	var encoded string
	base := big.NewInt(int64(len(encodeStd)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		encoded = string(encodeStd[mod.Int64()]) + encoded
	}
	return encoded
}

func ShowShorten(lite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			originalURL = "https://" + originalURL
		}

		shortenedURL := shortenUrl(originalURL)
		if err := db.StoreURL(lite, originalURL, shortenedURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]string{
			"OriginalURL": originalURL,
			"ShortURL":    shortenedURL,
		}

		tmpl, err := template.ParseFiles("internal/views/shorten.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Proxy(lite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:]
		if shortURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
		originalURL, err := db.GetOriginalURL(lite, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
	}
}
