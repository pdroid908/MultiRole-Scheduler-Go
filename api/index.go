package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"play/database"
	"play/middleware"
	"play/modules/auth"
	"play/redis"
	"play/user"
)

var app *gin.Engine

func init() {
	_ = godotenv.Load()

	pool, err := database.Connect()
	if err != nil {
		log.Printf("gagal konek ke database %v", err)
	}

	authData := &modules.Data{
		DB: pool,
	}

	userData := &user.DB{
		Database: pool,
	}

	err = redis.Connect()
	if err != nil {
		fmt.Println("gagal konek ke Redis")
	} else {
		log.Println("Berhasil terhubung ke Redis!")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://goojadwal.pages.dev", "http://localhost:5173", "https://goojadwal.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- ENDPOINT API ---
	r.POST("/regis", authData.Register)
	r.POST("/login", authData.Login)
	r.POST("/tes", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok": "ok",
		})
	})

	r.GET("/api/public/jadwal/:id", userData.MemberProfile)
	r.GET("/api/booking/jadwal/:username", userData.BookingProfile)
	r.POST("/api/booking/:username", userData.CreateBooking)

	// --- RUTE PROTECTED ---
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/api/download", func(c *gin.Context) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		}, userData.Profile)

		protected.POST("/logout", authData.Logout)
		protected.GET("/user", userData.Profile)
		protected.PUT("/user/accept/:id", userData.Accept)
		protected.DELETE("/user/jadwal/:id", userData.Delete)
		protected.PUT("/user/jadwal/keterangan/:id", userData.UpdateKeterangan)
		protected.PUT("/user/change-password", userData.UpdatePassword)
		protected.PUT("/user/change-username", userData.UpdateUsername)
	}

	app = r
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}