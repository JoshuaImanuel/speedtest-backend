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

// Handler untuk tes download (versi paling stabil untuk Vercel)
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    // --- PERUBAHAN DI SINI ---
	// Mengurangi ukuran total menjadi 5MB untuk memastikan stabilitas maksimum
	totalSize := 5 * 1024 * 1024      
	chunkSize := 256 * 1024           
	chunk := make([]byte, chunkSize)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", totalSize))

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for bytesSent := 0; bytesSent < totalSize; bytesSent += chunkSize {
		remaining := totalSize - bytesSent
		if remaining < chunkSize {
			chunk = make([]byte, remaining)
		}
		if _, err := rand.Read(chunk); err != nil {
			return
		}
		if _, err := w.Write(chunk); err != nil {
			return
		}
		flusher.Flush()
	}
}


// Handler untuk tes upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.Copy(io.Discard, r.Body)
	if err != nil {
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})
	handler := c.Handler(mux)
	handler.ServeHTTP(w, r)
}

