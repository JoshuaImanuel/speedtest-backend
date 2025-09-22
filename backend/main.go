package handler

import (
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

// Handler untuk tes download (versi ultra ringan)
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	totalSize := 15 * 1024 * 1024      // Total 15MB (ukuran yang aman)
	chunkSize := 256 * 1024           // Potongan 256KB
	
	// Buat satu potongan data statis yang akan dikirim berulang kali.
	// Ini sangat ringan dan hampir tidak menggunakan CPU.
	chunk := make([]byte, chunkSize)
	for i := range chunk {
		chunk[i] = '0' // Isi dengan karakter '0'
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", totalSize))

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for bytesSent := 0; bytesSent < totalSize; bytesSent += chunkSize {
		// Menentukan ukuran potongan terakhir
		remaining := totalSize - bytesSent
		if remaining < chunkSize {
			chunk = chunk[:remaining]
		}
		
		// Tulis potongan data statis
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

