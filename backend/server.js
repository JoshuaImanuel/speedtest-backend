// 1. Impor library yang kita butuhkan
// Fastify adalah framework web kita, sangat cepat dan efisien.
const fastify = require('fastify');
const crypto = require('crypto'); // Modul bawaan Node.js untuk membuat data acak

// 2. Buat instance server Fastify
const server = fastify();

// 3. Konfigurasi CORS (Cross-Origin Resource Sharing)
// INI SANGAT PENTING!
// Frontend kita berjalan di port yang berbeda (dari Live Server, misal: 5500)
// dengan backend kita (port 3000). Browser akan memblokir permintaan
// antar port yang berbeda kecuali backend secara eksplisit mengizinkannya.
// Inilah fungsi dari CORS.
server.register(require('@fastify/cors'), {
  origin: "*", // Mengizinkan semua origin (untuk pengembangan).
});

// 4. Endpoint untuk tes PING / LATENCY
// Frontend akan mengirim request ke URL ini dan mengukur waktu bolak-balik.
server.get('/ping', async (request, reply) => {
  reply.send({ status: 'ok' });
});

// 5. Endpoint untuk tes DOWNLOAD
// Kita siapkan sebuah "bongkahan" data acak sebesar 100MB di memori.
// Saat frontend meminta URL ini, kita akan kirim data ini.
const downloadData = crypto.randomBytes(100 * 1024 * 1024); // 100 MB
server.get('/download', async (request, reply) => {
  // Set header agar browser tahu ini adalah data untuk diunduh (bukan halaman web)
  reply.header('Content-Type', 'application/octet-stream');
  reply.send(downloadData);
});

// 6. Endpoint untuk tes UPLOAD
// Frontend akan mengirim data ke URL ini.
// Tugas server hanyalah menerima request tersebut dan merespons 'ok'.
// Kita tidak perlu melakukan apa-apa dengan data yang dikirim.
server.post('/upload', async (request, reply) => {
  // Cukup kirim status 'ok' untuk menandakan data telah diterima.
  reply.send({ status: 'received' });
});

// 7. Fungsi untuk menjalankan server
const start = async () => {
  try {
    // Server akan "mendengarkan" semua permintaan yang masuk ke port 3000
    await server.listen({ port: 3000 });
    console.log(`Server backend berjalan di http://localhost:3000`);
  } catch (err) {
    console.error(err);
    process.exit(1);
  }
};

// 8. Panggil fungsi untuk menjalankan server
start();

