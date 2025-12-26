package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"ziyad-test/database"
	"ziyad-test/model"
	"ziyad-test/utils"

	_ "github.com/lib/pq"
)

// handler

func BorrowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.BorrowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", "ZYD-ERR-400")
		return
	}

	// transaksi database
	// Kita perlu memastikan bahwa pengecekan stok, kuota, dan proses update dilakukan secara atomik.
	// Jika dua user meminjam buku terakhir di waktu yang sama atau race condition ,
	// database isolation level akan menjamin hanya satu yang berhasil.
	ctx := context.Background()
	tx, err := database.DB.BeginTx(ctx, nil)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Gagal memulai transaksi", "ZYD-ERR-500")
		return
	}

	// Defer Rollback Pendekatan ini dipilih agar jika terjadi error di tengah jalan
	// atau aplikasi panik, perubahan data tidak akan tersimpan secara tidak sengaja.
	defer tx.Rollback()

	// 1.Validasi Stok Buku dengan FOR UPDATE
	// Baris ini akan dikunci sehingga proses lain tidak bisa mengubah stok buku ini
	// sampai transaksi ini selesai. Ini krusial untuk mencegah "double booking" stok.
	var stock int
	err = tx.QueryRowContext(ctx, "SELECT stock FROM books WHERE id = $1 FOR UPDATE", req.BookID).Scan(&stock)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Buku tidak ditemukan", "ZYD-ERR-002")
		return
	}

	if stock <= 0 {
		utils.SendError(w, http.StatusConflict, "Stok buku habis", "ZYD-ERR-001")
		return
	}

	// 2.Validasi Kuota Member
	var quota int
	err = tx.QueryRowContext(ctx, "SELECT quota FROM members WHERE id = $1 FOR UPDATE", req.MemberID).Scan(&quota)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Member tidak ditemukan", "ZYD-ERR-003")
		return
	}

	if quota <= 0 {
		utils.SendError(w, http.StatusConflict, "Kuota peminjaman member sudah penuh", "ZYD-ERR-004")
		return
	}

	// 3.Eksekusi Perubahan Data
	_, err = tx.ExecContext(ctx, "UPDATE books SET stock = stock - 1 WHERE id = $1", req.BookID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Gagal update stok", "ZYD-ERR-500")
		return
	}

	_, err = tx.ExecContext(ctx, "UPDATE members SET quota = quota - 1 WHERE id = $1", req.MemberID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Gagal update kuota", "ZYD-ERR-500")
		return
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (member_id, book_id, borrow_date) VALUES ($1, $2, NOW())", req.MemberID, req.BookID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Gagal mencatat transaksi", "ZYD-ERR-500")
		return
	}

	// Commit Transaksi
	if err := tx.Commit(); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Gagal commit transaksi", "ZYD-ERR-500")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Peminjaman berhasil",
		"status":  "success",
	})
}
