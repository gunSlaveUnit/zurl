package main

import (
	"crypto/sha256"
    "encoding/hex"
    "net/http"
)

var cache = make(map[string]string)

func redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("url")
	if longURL, exists := cache[shortURL]; exists {
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	} 
}

func short(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("longURL")

	hash := sha256.New()
    hash.Write([]byte(longURL))
    hashed := hash.Sum(nil)

	shortURL := hex.EncodeToString(hashed)[:6]

	cache[shortURL] = longURL

	w.Write([]byte(shortURL))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/redirect/{url}", redirect)
	mux.HandleFunc("/api/v1/short", short)

	http.ListenAndServe(":8080", mux)
}