// 1. Impor library yang kita butuhkan
const fastify = require('fastify');
const crypto = require('crypto'); 

// 2. Buat instance server Fastify
const server = fastify();

// 3. Konfigurasi CORS (Cross-Origin Resource Sharing)
server.register(require('@fastify/cors'), {
  origin: "*", 
});

// 4. Endpoint untuk tes PING / LATENCY
server.get('/ping', async (request, reply) => {
  reply.send({ status: 'ok' });
});

// 5. Endpoint untuk tes DOWNLOAD
const downloadData = crypto.randomBytes(100 * 1024 * 1024); // 100 MB
server.get('/download', async (request, reply) => {
  // Set header agar browser tahu ini adalah data untuk diunduh (bukan halaman web)
  reply.header('Content-Type', 'application/octet-stream');
  reply.send(downloadData);
});

// 6. Endpoint untuk tes UPLOAD
server.post('/upload', async (request, reply) => {
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

