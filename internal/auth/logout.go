package auth

import (
	"net/http"
	

	"github.com/gin-gonic/gin"
)

func (d *Data) Logout(c *gin.Context) {
	

	// 2. Hapus Cookie Token di Browser
	// Mengisi MaxAge = -1 dan value = "" akan memaksa browser langsung membuang cookie ini
	c.SetCookie(
		"token",
		"",
		-1, // MaxAge -1 = Hapus cookie seketika
		"/",
		"",
		false,
		true,
	)

	// 4. Return respon sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil logout",
	})
}