package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingRequest struct {
	Tanggal    string  `json:"tanggal" binding:"required"`
	JamMulai   string  `json:"jam_mulai" binding:"required"`
	JamSelesai string  `json:"jam_selesai" binding:"required"`
	Keterangan *string `json:"keterangan"`
}

func (d *DB) CreateBooking(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Username target tidak ditemukan"})
		return
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Data booking tidak lengkap atau format tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 1. Cari user_id pemilik jadwal (owner) berdasarkan username
	var ownerID string
	err := d.Database.QueryRow(ctx, "SELECT id::text FROM pengguna WHERE username = $1", username).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "Owner jadwal tidak ditemukan"})
		return
	}

	// Format jam memastikan HH:MM:SS
	jamMulai := req.JamMulai
	if len(jamMulai) == 5 {
		jamMulai += ":00"
	}
	jamSelesai := req.JamSelesai
	if len(jamSelesai) == 5 {
		jamSelesai += ":00"
	}

	// 2. Cek Bentrok Jadwal (Hanya mengecek jadwal yang SUDAH DI-ACC / is_confirmed = true)
	var count int
	checkOverlapQuery := `
		SELECT COUNT(*) 
		FROM jadwal 
		WHERE user_id = $1 
		  AND tanggal = $2::date
		  AND is_confirmed = true
		  AND jam_mulai < $3::time 
		  AND jam_selesai > $4::time
	`
	err = d.Database.QueryRow(ctx, checkOverlapQuery, ownerID, req.Tanggal, jamSelesai, jamMulai).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal mengecek ketersediaan jadwal: " + err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"err": "Jam tersebut sudah dibooking & di-ACC orang lain! Silakan pilih jam lain.",
		})
		return
	}

	// 3. Insert ke database jadwal
	insertQuery := `
		INSERT INTO jadwal (user_id, tanggal, jam_mulai, jam_selesai, keterangan, is_confirmed)
		VALUES ($1, $2::date, $3::time, $4::time, $5, false)
	`
	_, err = d.Database.Exec(ctx, insertQuery, ownerID, req.Tanggal, jamMulai, jamSelesai, req.Keterangan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal menyimpan jadwal booking: " + err.Error()})
		return
	}

	// Clear cache agar data publik/member langsung terbarui
	d.InvalidateUserCache(ownerID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Booking berhasil diajukan! Menunggu konfirmasi pemilik.",
	})
}