package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
)

var (
	urlsMap      = make(map[string]string)
	urlsMapMutex sync.RWMutex
)

func generateShortID() (string, error) {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	longURL := r.URL.Query().Get("url")
	if longURL == "" {
		http.Error(w, "Missing url query parameter", http.StatusBadRequest)
		return
	}

	shortID, err := generateShortID()

	if err != nil {
		http.Error(w, "Error generating short ID", http.StatusInternalServerError)
		return
	}

	urlsMapMutex.Lock()
	urlsMap[shortID] = longURL
	urlsMapMutex.Unlock()

	shortURL := fmt.Sprintf("http://%s/%s", r.Host, shortID)
	fmt.Fprintf(w, "Short URL : %s\n", shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[1:]

	urlsMapMutex.RLock()
	longURL, ok := urlsMap[shortID]
	urlsMapMutex.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)

}

func main() {
	http.HandleFunc("/shorten", shortURLHandler)
	http.HandleFunc("/", redirectHandler)

	port := "8080"
	fmt.Println("Starting the Url shortener...", port)
	http.ListenAndServe(":"+port, nil)
}
