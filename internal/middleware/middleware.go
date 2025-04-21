package middleware

import (
	"log"

	"github.com/alexey-shedrin/avito-test-task/internal/utils/token"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	jwtSecret           = "secretKey"
)

func Auth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.SetPrefix("middleware.Middleware")
		jwt := c.GetHeader(AuthorizationHeader)
		claims, err := token.ValidateJWT(jwt, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})

			return
		}
		role := claims.Role
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "unauthorized"})
	}
}
