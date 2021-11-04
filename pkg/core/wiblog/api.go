package wiblog

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @title APP Demo API
// @version 1.0
// @description This is a sample server celler server.

// AuthFilter auth filter
func AuthFilter(c *gin.Context) {
	if !IsLogined(c) {
		c.Abort()
		c.Status(http.StatusUnauthorized)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
}

// SetLogin login user
func SetLogin(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Set("username", username)
	_ = session.Save()
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
