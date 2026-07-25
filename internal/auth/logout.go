package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Data) Logout(c *gin.Context) {
	// 1. Ambil user_id / username jika pernah disimpan di cache (opsional)
	userID, exists := c.Get("user_id")
	if exists && d.Cache != nil {
		if idStr, ok := userID.(string); ok {
			// Hapus data session di go-cache jika ada
			d.Cache.Delete("user_session:" + idStr)
		}
	}

	// 2. Hapus Cookie Token di Browser
	// Mengisi MaxAge = -1 dan value = "" akan memaksa browser langsung membuang cookie ini
	c.SetCookie(
		"token",
		"",
		-1, // MaxAge -1 = Hapus cookie seketika
		"/",
		"localhost",
		false,
		true,
	)

	// 3. Tambahkan Header agar Browser TIDAK Menimpan Cache Halaman
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 4. Return respon sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil logout",
	})
}