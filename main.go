package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type URLShortener struct {
	Index       string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

var SQL = make(map[string]URLShortener)

func Generate(originalURL string) string {
	val := md5.New()
	val.Write([]byte(originalURL))
	return hex.EncodeToString(val.Sum(nil))[:8]
}

func createURL(originalURL string) string {
	shortURL := Generate(originalURL)
	SQL[shortURL] = URLShortener{
		Index:       shortURL,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	fmt.Println("Shortened URL:", "http://localhost:8080/redirect/"+shortURL)
	return shortURL
}

func getURL(identity string) (URLShortener, error) {
	Address, ok := SQL[identity]
	if !ok {
		return URLShortener{}, errors.New("Error: 404 Not Found")
	}
	return Address, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to URL Shortener Service")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var info struct {
		OriginalURL string `json:"original_url"`
	}
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL_ := createURL(info.OriginalURL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	identity := r.URL.Path[len("/redirect/"):]
	Address, err := getURL(identity)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, Address.OriginalURL, http.StatusFound)
}

func main() {

	url := ""
	fmt.Println("Enter the URL to be shortened: ")
	fmt.Scanln(&url)

	shortURL := createURL(url)
	fmt.Println("Shortened URL:", "http://localhost:8080/redirect/"+shortURL)

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", RedirectHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error on starting server:", err)
	}
}
