package main

import (
	"context"
	"crypto/sha256"
    "encoding/hex"
    "net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	cache *redis.Client
)

func redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("url")
	longURL, err := cache.Get(ctx, shortURL).Result()
	if err == nil {
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

	cache.Set(ctx, shortURL, longURL, 0)

	w.Write([]byte(shortURL))
}

func main() {
	godotenv.Load()

	cache_db, _ := strconv.Atoi(os.Getenv("CACHE_DB"))
	cache = redis.NewClient(&redis.Options{
        Addr: os.Getenv("CACHE_URL"),
		DB: cache_db,
        Password: os.Getenv("CACHE_PASSWORD"),
    })

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/redirect/{url}", redirect)
	mux.HandleFunc("/api/v1/short", short)

	http.ListenAndServe(":8080", mux)
}