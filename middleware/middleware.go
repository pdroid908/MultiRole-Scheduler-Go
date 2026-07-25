package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": " invalid token"})
			c.Abort()
			return
		}

		jtwKey := []byte(os.Getenv("JWT_SECRET"))

		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokens, claims, func(token *jwt.Token) (interface{}, error) {
			return jtwKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "tidak valid token"})
			c.Abort()
			return
		}

		if Id, ok := (*claims)["user_id"].(string); ok {
			c.Set("user_id", Id)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"err": " user id invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
