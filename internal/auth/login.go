package auth

import (
	"context"
	"net/http"
	"os"
	// "play/redis"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email string `json:"email"`
}

func (d*Data) Login(c *gin.Context){
	var req login
	jwtKey:=[]byte(os.Getenv("JWT_SECRET"))

	if err:= c.ShouldBindJSON(&req); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":"invalid data"})
		return
	}
	if req.Email!=""{
		c.JSON(http.StatusOK,gin.H{"message":"sukses Login"})
		return
	}
	ctx,cancel:= context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()

	var passDB string
	var user_id string
	cari:= `select password_hash, id from pengguna where username= $1`
	if err:= d.DB.QueryRow(ctx,cari,req.Username).Scan(&passDB,&user_id); err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"err":"gagal ambil di db"}) 
		return
	}

	err:= bcrypt.CompareHashAndPassword([]byte(passDB),[]byte(req.Password))
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":"invalid data"})
		return
	}

	// limit:= "Login_Limit:"+ c.ClientIP()
	// boleh, err:= redis.RateLimit(limit)
	// if err!=nil{
	// 	c.JSON(http.StatusInternalServerError,gin.H{"err":"redis err"})
	// 	return
	// }
	// if !boleh{
	// 	c.JSON(http.StatusTooManyRequests,gin.H{"err":"tunggu 5 menit"})
	// 	return
	// }

	expJwt:= time.Now().Add(8*time.Hour).Unix()

	Claims:= &jwt.MapClaims{
		"user_id": user_id,
		"username":req.Username,
		"exp": expJwt,
	}

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,Claims)

	tokens, err:= token.SignedString(jwtKey)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"err":"servr err"})
		return
	}

	c.SetCookie(
		"token",
		tokens,
		int(8*time.Hour/time.Second),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK,gin.H{"message":"sukses login"})

}