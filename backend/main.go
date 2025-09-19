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
	// Total ukuran data yang akan dikirim
	totalSize := 25 * 1024 * 1024 
	
	// Ukuran setiap potongan data (misal: 1MB)
	chunkSize := 1 * 1024 * 1024 
	chunk := make([]byte, chunkSize)

	// Set header terlebih dahulu
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", totalSize))

	// Kirim data dalam potongan-potongan kecil
	for bytesSent := 0; bytesSent < totalSize; bytesSent += chunkSize {
		// Isi potongan dengan data acak
		rand.Read(chunk)
		// Tulis potongan ke response
		_, err := w.Write(chunk)
		if err != nil {
			// Hentikan jika koneksi dari klien terputus
			return
		}
		// Flush data untuk memastikan data terkirim langsung (penting)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
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

