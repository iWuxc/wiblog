package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// UserMiddleware 用户cookie标记
func UserMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Cookie("u")
		if err != nil || cookie != "" {
			u1 := uuid.Must(uuid.NewV4()).String()
			context.SetCookie("u", u1, 86400*730, "/", "", true, true)
		}
	}
}
