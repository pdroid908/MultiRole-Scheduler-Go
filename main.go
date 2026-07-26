package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"play/database"
	"play/internal/auth"
	"play/middleware"
	"play/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

// func SecurityHeaders() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("X-Frame-Options", "DENY")
// 		c.Header("X-Content-Type-Options", "nosniff")
// 		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
// 		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
// 		c.Header("X-XSS-Protection", "0")

// 		// PERBAIKAN CSP: Diizinkan CDN Tailwind, FontAwesome, dan Google Fonts agar UI tidak hancur
// 		cspPolicy := "default-src 'self'; " +
// 			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.tailwindcss.com; " +
// 			"style-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com https://fonts.googleapis.com; " +
// 			"font-src 'self' https://cdnjs.cloudflare.com https://fonts.gstatic.com; " +
// 			"img-src 'self' data: https:; " +
// 			"connect-src 'self'; " +
// 			"object-src 'none'; " +
// 			"frame-ancestors 'none'; " +
// 			"base-uri 'self'; " +
// 			"form-action 'self'"

// 		c.Header("Content-Security-Policy", cspPolicy)
// 		c.Next()
// 	}
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("gagal load env")
	}

	pool, err := database.Connect()
	if err != nil {
		log.Fatalf("gagal konek ke database %v", err)
	}
	defer pool.Close()

	c := cache.New(24*time.Hour, 10*time.Minute)

	auth := &auth.Data{
		DB:    pool,
		Cache: c,
	}

	user := &user.DB{
		Database: pool,
		Cache:    c,
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Sesuaikan port frontend React kamu
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	
	

	r.POST("/regis", auth.Register)
	r.POST("/login", auth.Login)
	
	
	r.GET("/api/public/jadwal/:id", user.MemberProfile)
	r.GET("/api/booking/jadwal/:username", user.BookingProfile)


	// Endpoint untuk booking dari publik/tamu tanpa login:
	r.POST("/api/booking/:username", user.CreateBooking)

	// --- 2. RUTE YANG BUTUH LOGIN (Masuk ke Group Protected) ---
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		

		protected.GET("/download", func(c *gin.Context) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		}, user.Profile)
        

		// Endpoint Logout
		protected.POST("/logout", auth.Logout)

		protected.GET("/user", user.Profile)
		protected.PUT("/user/accept/:id", user.Accept)
		protected.DELETE("/user/jadwal/:id", user.Delete)
		protected.PUT("/user/jadwal/keterangan/:id", user.UpdateKeterangan)

		// --- ENDPOINT BARU GANTI PASSWORD & USERNAME ---
		protected.PUT("/user/change-password", user.UpdatePassword)
		protected.PUT("/user/change-username", user.UpdateUsername)
	}


	srv := &http.Server{
    Addr: ":8080",
    Handler: r,

    ReadTimeout: 10 * time.Second,
    ReadHeaderTimeout: 5 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout: 60 * time.Second,
}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server gagal jalan %v", err)
		}
	}()

	log.Print("server jalan di 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("mematikan server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server down gagal shutdown", err)
	}

	log.Print("server berhenti dengan aman")
}
