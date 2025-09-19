package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

// Handler untuk tes ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

// Handler untuk tes download
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Buat data acak sebesar 100MB
	dataSize := 100 * 1024 * 1024
	randomData := make([]byte, dataSize)
	rand.Read(randomData)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", dataSize))
	w.Write(randomData)
}

// Handler untuk tes upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Baca dan buang seluruh body dari request untuk mengukur upload
	// io.Copy dengan io.Discard adalah cara paling efisien
	_, err := io.Copy(io.Discard, r.Body)
	if err != nil {
		http.Error(w, "Failed to read upload data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Fprintf(w, `{"status":"ok"}`)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/upload", uploadHandler)

	// Konfigurasi CORS agar frontend bisa mengakses
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Izinkan semua
		AllowedMethods: []string{"GET", "POST"},
	})
	handler := c.Handler(mux)

	port := "8080"
	fmt.Printf("Go server running on port %s\n", port)
	
	// Gunakan server kustom untuk mengatur timeout
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
