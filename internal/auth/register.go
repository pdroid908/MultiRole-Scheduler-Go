package auth

import (

	"context"

	"net/http"

	"play/redis"

	// "play/redis"
	"time"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	
)

type register struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	// Honeypot field: tidak boleh diisi oleh manusia
	Email string `json:"email"`
}

func (d *Data) Register(c *gin.Context){
	var req register
	limit := "Regis_Limit:" + c.ClientIP()
	boleh, err:= redis.RateLimit(limit)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"err":"invalid err"})
		return
	}
	if !boleh{
		c.JSON(http.StatusTooManyRequests,gin.H{"err":"tunggu 5 ment"})
		return
	}
	if req.Email!=""{
		c.JSON(http.StatusOK,gin.H{"message":"regis sukses"})
		return
	}
	if err:= c.ShouldBindJSON(&req); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":"data invalid"})
		return
	}

	ctx, cancel:= context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()

	var udhAda bool
	cari:= `select exists( select 1 from pengguna where username =$1)`
	err = d.DB.QueryRow(ctx,cari,req.Username).Scan(&udhAda)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"err":"gagal cari db"})
		return
	}
	if udhAda{
		c.JSON(http.StatusBadRequest,gin.H{"err":"invalid akun"})
		return
	}

	hashPass,err:= bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"err":"gagal genr pass"})
		return
	}

	NewAkun:= `insert into pengguna (username, password_hash) values ($1, $2)`
	_,err = d.DB.Exec(ctx,NewAkun,req.Username,string(hashPass))
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"err":"gagal masuk insert"})
		return
	}

	c.JSON(http.StatusCreated,gin.H{"message":"sukses regis"})

}

