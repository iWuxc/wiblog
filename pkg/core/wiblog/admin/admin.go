package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"wiblog/pkg/cache"
	"wiblog/pkg/core/wiblog"
	"wiblog/tools"
)

func RegisterRoutes(e *gin.Engine) {
	e.POST("/admin/login", handleAcctLogin)
}

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


}
