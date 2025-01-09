package main

import (
	"context"
	"crypto/sha256"
    "encoding/hex"
	"html/template"
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

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		longURL := r.FormValue("longURL")

		hash := sha256.New()
		hash.Write([]byte(longURL))
		hashed := hash.Sum(nil)

		shortURL := hex.EncodeToString(hashed)[:6]

		cache.Set(ctx, shortURL, longURL, 0)

		t, _ := template.ParseFiles("templates/index.html")
		data := map[string]interface{}{
			"ShortURL": "http://localhost:8080/" + shortURL,
		}

		t.Execute(w, data)
	}
	
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, nil)
	}
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

	mux.HandleFunc("/", home)
	mux.HandleFunc("/api/v1/redirect/{url}", redirect)

	http.ListenAndServe(":8080", mux)
}