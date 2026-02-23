package middleware

import (
	"net/http"
	"restfulapi/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authheader := ctx.GetHeader("Authorization")

		if authheader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization",
			})

			ctx.Abort()
			return
		}

		parts := strings.Split(authheader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization",
			})
			ctx.Abort()
			return
		}

		tokenstring := parts[1]
		claims, err := utils.ValidateJWT(tokenstring)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token",
			})

			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}
