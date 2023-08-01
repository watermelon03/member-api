package middlewere

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authen() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		auth := strings.Split(tokenString, " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "forbidden",
				"message": "Authorization failed",
			})
			return
		}

		mySigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("adminID", claims["adminID"])
			c.Set("userID", claims["userID"])
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "forbidden",
				"message": err.Error(),
			})
		}
		c.Next()
	}
}
