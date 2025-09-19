package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/cors"
)

// Handler untuk tes ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

// Handler untuk tes download
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	dataSize := 100 * 1024 * 1024
	randomData := make([]byte, dataSize)
	rand.Read(randomData)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", dataSize))
	w.Write(randomData)
}

// Handler untuk tes upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.Copy(io.Discard, r.Body)
	if err != nil {
		http.Error(w, "Failed to read upload data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Fprintf(w, `{"status":"ok"}`)
}

// Handler utama yang akan dipanggil oleh Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/upload", uploadHandler)

	// Konfigurasi CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})
	handler := c.Handler(mux)

	// Jalankan handler yang sudah dibungkus CORS
	handler.ServeHTTP(w, r)
}


