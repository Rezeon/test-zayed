package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// database handler

var DB *sql.DB

func InitDB() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, menggunakan environment variables dari sistem")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Melakukan retry koneksi karena DB mungkin belum siap saat Docker naik
	for i := 0; i < 5; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil && DB.Ping() == nil {
			break
		}
		log.Printf("Menunggu database siap (percobaan %d/5)...", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Membuat skema tabel sederhana
	schema := `
	CREATE TABLE IF NOT EXISTS books (id SERIAL PRIMARY KEY, title TEXT, stock INT);
	CREATE TABLE IF NOT EXISTS members (id SERIAL PRIMARY KEY, name TEXT, quota INT);
	CREATE TABLE IF NOT EXISTS transactions (id SERIAL PRIMARY KEY, member_id INT, book_id INT, borrow_date TIMESTAMP);
	
	INSERT INTO books (title, stock) SELECT 'Belajar Go', 1 WHERE NOT EXISTS (SELECT 1 FROM books WHERE id = 1);
	INSERT INTO members (name, quota) SELECT 'Ziyad', 2 WHERE NOT EXISTS (SELECT 1 FROM members WHERE id = 1);
	`
	DB.Exec(schema)
}
