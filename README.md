Sebelum memulai, pastikan Anda telah menginstal:

Docker

Docker Compose

Arsitektur Stack

Backend: Go (Golang)

Database: PostgreSQL 15 (Alpine)

Orchestration: Docker Compose

Konfigurasi Variabel Lingkungan

Aplikasi menggunakan variabel lingkungan berikut yang dikonfigurasi di dalam docker-compose.yaml:

Variabel

Deskripsi

Nilai Default

DB_HOST

Host database (nama service docker)

db

DB_PORT

Port internal database

5432

DB_USER

Username PostgreSQL

postgres

DB_PASSWORD

Password PostgreSQL

12345678

DB_NAME

Nama Database

ziyadtest

Cara Menjalankan Aplikasi

0. clone project ini 

git clone https://github.com/Rezeon/test-zayed.git

1. Build dan Jalankan Kontainer

Gunakan perintah berikut untuk membangun image dan menjalankan seluruh layanan:

docker-compose up --build


2. Menjalankan di Background 

Jika Anda ingin menjalankan kontainer di latar belakang:

docker-compose up -d


3. Menghentikan Layanan

Untuk menghentikan dan menghapus kontainer:

docker-compose down
