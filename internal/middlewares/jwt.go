package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func HasAuthorization(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token requerido"})
			return
		}

		tk, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
			return
		}

		if !tk.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
			return
		}

		userID, ok := tk.Claims.(jwt.MapClaims)["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
