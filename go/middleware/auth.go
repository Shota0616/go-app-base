package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go-app-base/auth"
    "log"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        log.Println(token)

        // リクエストにトークンが載っていなかったらエラーを返す
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
            c.Abort()
            return
        }

        // リクエストのトークンがvalidate通らなかったらエラーを返す
        claims, err := auth.ValidateJWT(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("id", claims.ID)
        c.Next()
    }
}
