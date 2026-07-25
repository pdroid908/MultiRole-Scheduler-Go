package user

import (
	"context"
	"net/http"
	"play/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
)

type DB struct {
	Database *pgxpool.Pool
	Cache    *cache.Cache
}

func NewDBService(pool *pgxpool.Pool) *DB {
	// Cache disimpan default 12 jam, cleanup item mati tiap 10 menit
	c := cache.New(12*time.Hour, 10*time.Minute)
	return &DB{
		Database: pool,
		Cache:    c,
	}
}

func (d *DB) InvalidateUserCache(userID string) {
	d.Cache.Delete("public_schedule:" + userID)
	d.Cache.Delete("admin_profile:" + userID)

	// Hapus juga cache booking berdasarkan username
	var username string
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := d.Database.QueryRow(ctx, "SELECT username FROM pengguna WHERE id = $1", userID).Scan(&username)
	if err == nil && username != "" {
		d.Cache.Delete("booking_schedule:" + username)
	}
}

// -------------------------------------------------------------
// 1. ADMIN PROFILE (TEMPAT MEMASAK & MENULIS CACHE PERTAMA KALI)
// -------------------------------------------------------------
func (d *DB) Profile(c *gin.Context) {
	userID, ada := c.Get("user_id")
	if !ada {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid id"})
		return
	}
	ID, valid := userID.(string)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid id ya"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var username string
	cari := `select username from pengguna where id=$1`
	err := d.Database.QueryRow(ctx, cari, ID).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"er": "server gagal ambil nama"})
		return
	}

	query := `
		SELECT 
			id::text, 
			user_id::text, 
			tanggal::text, 
			jam_mulai::text, 
			jam_selesai::text, 
			keterangan, 
			is_confirmed 
		FROM jadwal 
		WHERE user_id = $1 
		ORDER BY tanggal ASC, jam_mulai ASC
	`

	rows, err := d.Database.Query(ctx, query, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "gagal ambil database: " + err.Error()})
		return
	}
	defer rows.Close()

	var daftarJadwal []models.Jadwal = []models.Jadwal{}
	var daftarBooking []BookingJadwal = []BookingJadwal{}

	for rows.Next() {
		var j models.Jadwal
		err := rows.Scan(
			&j.Id,
			&j.UserId,
			&j.Tanggal,
			&j.JamMulai,
			&j.JamSelesai,
			&j.Keterangan,
			&j.IsConfirmed,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "gagal membongkar data jadwal"})
			return
		}

		daftarJadwal = append(daftarJadwal, j)

		// Konversi ke format BookingJadwal (Tanpa keterangan untuk privasi member)
		daftarBooking = append(daftarBooking, BookingJadwal{
			Id:          j.Id,
			UserId:      j.UserId,
			Tanggal:     j.Tanggal,
			JamMulai:    j.JamMulai,
			JamSelesai:  j.JamSelesai,
			IsConfirmed: j.IsConfirmed,
		})
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "error saat membaca baris data"})
		return
	}

	// -------------------------------------------------------------
	// LANGSUNG MASAK & WRITE KEDUA CACHE SAAT ADMIN AKSES PROFIL!
	// -------------------------------------------------------------
	d.Cache.Set("public_schedule:"+ID, daftarJadwal, cache.NoExpiration)
	d.Cache.Set("booking_schedule:"+username, daftarBooking, cache.NoExpiration)

	// Kirim data ke Admin
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"id":       ID,
		"username": username,
		"data":     daftarJadwal,
	})
}

// -------------------------------------------------------------
// 2. MEMBER PROFILE (STRICTLY BACA CACHE ONLY - NO QUERY DB)
// -------------------------------------------------------------
func (d *DB) MemberProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "id tidak boleh kosong"})
		return
	}

	cacheKey := "public_schedule:" + userID

	// BACA DARI CACHE SAJA
	if cachedData, found := d.Cache.Get(cacheKey); found {
		if jadwal, ok := cachedData.([]models.Jadwal); ok {
			c.JSON(http.StatusOK, jadwal)
			return
		}
	}

	// JIKA CACHE BELUM DIASAK ADMIN -> RETURNING EMPTY ARRAY (TIDAK GET DB)
	c.JSON(http.StatusOK, []models.Jadwal{})
}

type BookingJadwal struct {
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	Tanggal     string `json:"tanggal"`
	JamMulai    string `json:"jam_mulai"`
	JamSelesai  string `json:"jam_selesai"`
	IsConfirmed bool   `json:"is_confirmed"`
}

// -------------------------------------------------------------
// 3. BOOKING PROFILE MEMBER (STRICTLY BACA CACHE ONLY - NO QUERY DB)
// -------------------------------------------------------------
func (d *DB) BookingProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "username tidak boleh kosong"})
		return
	}

	cacheKey := "booking_schedule:" + username

	// BACA DARI CACHE SAJA
	if cachedData, found := d.Cache.Get(cacheKey); found {
		if jadwal, ok := cachedData.([]BookingJadwal); ok {
			c.JSON(http.StatusOK, jadwal)
			return
		}
	}

	// JIKA CACHE BELUM DIMASAK ADMIN -> RETURNING EMPTY ARRAY (TIDAK GET DB)
	c.JSON(http.StatusOK, []BookingJadwal{})
}