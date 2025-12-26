package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"ziyad-test/database"
	"ziyad-test/handler"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	godotenv.Load()
	database.InitDB()
	defer database.DB.Close()

	http.HandleFunc("/borrow", handler.BorrowHandler)

	fmt.Println("Server Library Transaction berjalan di :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
