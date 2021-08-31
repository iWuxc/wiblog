package wiblog

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetLogin login user
func SetLogin(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()
}

// SetLogout logout user
func SetLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("username")
	_ = session.Save()
}

// IsLogined account logined
func IsLogined(c *gin.Context) bool {
	return GetUsername(c) != ""
}

// GetUsername get logined account
func GetUsername(c *gin.Context) string {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		return ""
	}
	return username.(string)
}
