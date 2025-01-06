package main

import (
	"crypto/sha256"
    "encoding/hex"
    "net/http"
)

var cache = make(map[string]string)

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

	mux.HandleFunc("/api/v1/short", short)

	http.ListenAndServe(":8080", mux)
}