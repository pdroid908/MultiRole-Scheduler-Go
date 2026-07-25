package user

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func(d *DB) Accept(c *gin.Context){
	UsrID, ada := c.Get("user_id")
	if !ada{
		c.JSON(http.StatusUnauthorized,gin.H{"err":"invalid id"})
		return
	}
	ID, valid:= UsrID.(string)
	if !valid{
		c.JSON(http.StatusUnauthorized,gin.H{"err":"invalid id"})
		return
	}
	ctx, cancel:= context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()

	jadwalID := c.Param("id")
	if jadwalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "id jadwal tidak boleh kosong"})
		return
	}

	query := `
		UPDATE jadwal 
		SET is_confirmed = true 
		WHERE id = $1 AND user_id = $2
	`
	result, err := d.Database.Exec(ctx, query, jadwalID, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "gagal update database: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"err": "jadwal tidak ditemukan atau kamu tidak memiliki akses"})
		return
	}
	d.InvalidateUserCache(ID)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Jadwal berhasil disetujui!",
	})
}

func (d *DB) Delete(c *gin.Context) {
	UsrID, ada := c.Get("user_id")
	if !ada {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid id"})
		return
	}
	ID, valid := UsrID.(string)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid id"})
		return
	}

	jadwalID := c.Param("id")
	if jadwalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "id jadwal tidak boleh kosong"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM jadwal WHERE id = $1 AND user_id = $2`
	result, err := d.Database.Exec(ctx, query, jadwalID, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "gagal menghapus data dari database: " + err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"err": "jadwal tidak ditemukan atau kamu tidak memiliki akses"})
		return
	}
	d.InvalidateUserCache(ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Jadwal berhasil dihapus!",
	})
}

type UpdateKeteranganRequest struct {
	Keterangan string `json:"keterangan"`
}

// 2. FUNC UPDATE KETERANGAN: Hanya memperbarui teks keterangan
func (d *DB) UpdateKeterangan(c *gin.Context) {
	UsrID, ada := c.Get("user_id")
	if !ada {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid id"})
		return
	}
	ID, valid := UsrID.(string)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid id"})
		return
	}

	jadwalID := c.Param("id")
	if jadwalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "id jadwal tidak boleh kosong"})
		return
	}

	var req UpdateKeteranganRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "format request tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE jadwal 
		SET keterangan = $1 
		WHERE id = $2 AND user_id = $3
	`
	result, err := d.Database.Exec(ctx, query, req.Keterangan, jadwalID, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "gagal memperbarui keterangan: " + err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"err": "jadwal tidak ditemukan atau kamu tidak memiliki akses"})
		return
	}
	d.InvalidateUserCache(ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Keterangan berhasil diperbarui!",
	})
}



// Struct Request untuk Ganti Password
type UpdatePasswordReq struct {
	PasswordLama string `json:"password_lama" binding:"required"`
	PasswordBaru string `json:"password_baru" binding:"required,min=6"`
}

// Struct Request untuk Ganti Username
type UpdateUsernameReq struct {
	UsernameBaru string `json:"username_baru" binding:"required,min=3,max=50"`
}

// 1. Fungsi Ganti Password
func (d *DB) UpdatePassword(c *gin.Context) {
	userID, ada := c.Get("user_id")
	if !ada {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized"})
		return
	}
	ID := userID.(string)

	var req UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Input tidak valid, minimal password 6 karakter"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Ambil password hash lama dari DB
	var currentHash string
	err := d.Database.QueryRow(ctx, "SELECT password_hash FROM pengguna WHERE id = $1", ID).Scan(&currentHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal mengambil data user"})
		return
	}

	// Cek apakah password lama sesuai
	err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.PasswordLama))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Password lama salah!"})
		return
	}

	// Hash password baru
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal memproses password baru"})
		return
	}

	// Update password ke DB
	_, err = d.Database.Exec(ctx, "UPDATE pengguna SET password_hash = $1 WHERE id = $2", string(newHash), ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal memperbarui password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password berhasil diubah!"})
}

// 2. Fungsi Ganti Username
// 2. Fungsi Ganti Username (Versi Lebih Cepat & Ringkas)
func (d *DB) UpdateUsername(c *gin.Context) {
	userID, ada := c.Get("user_id")
	if !ada {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized"})
		return
	}
	ID, valid := userID.(string)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized"})
		return
	}

	var req UpdateUsernameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Username minimal 3 karakter"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Langsung UPDATE ke DB tanpa SELECT EXISTS dulu
	_, err := d.Database.Exec(ctx, "UPDATE pengguna SET username = $1 WHERE id = $2", req.UsernameBaru, ID)
	if err != nil {
		// Menangkap error jika username sudah dipakai orang lain (Constraint Unique PostgreSQL)
		// Misal pesan error dari pgx mengandung 'unique' atau '23505'
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"err": "Username sudah digunakan oleh akun lain"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal memperbarui username: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Username berhasil diperbarui!"})
}