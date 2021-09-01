package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"wiblog/pkg/cache"
	"wiblog/pkg/core/wiblog"
	"wiblog/tools"
)

// RegisterRoutes register routes
func RegisterRoutes(e *gin.Engine) {
	e.POST("/admin/login", handleAcctLogin)
}

// RegisterRoutesAuthz register routes
func RegisterRoutesAuthz(group *gin.IRoutes) {

}

// handleAcctLogin login passport
func handleAcctLogin(c *gin.Context) {
	user := c.PostForm("username")
	passwd := c.PostForm("password")
	if user == "" || passwd == "" {
		logrus.Warnf("参数错误: %s, %s", user, passwd)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	if cache.Wi.Account.Username != user ||
		cache.Wi.Account.Password != tools.EncryptPassword(user, passwd) {
		logrus.Warnf("账号或者密码错误: %s, %s", user, passwd)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	//验证通过
	wiblog.SetLogin(c, user)
	//更新登录状态
	cache.Wi.Account.LoginIP = c.ClientIP()
	cache.Wi.Account.LoginAt = time.Now()
	_ = cache.Wi.Store.UpdateAccount(context.Background(), user, map[string]interface{}{
		"login_ip": cache.Wi.Account.LoginIP,
		"login_at": cache.Wi.Account.LoginAt,
	})
	c.Redirect(http.StatusFound, "/admin/profile")
}
