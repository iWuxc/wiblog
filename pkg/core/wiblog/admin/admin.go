package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"wiblog/pkg/cache"
	"wiblog/pkg/core/wiblog"
	"wiblog/tools"
)

var (
	NoticeSuccess = "success"
	NoticeWarning = "warning"
	NoticeError   = "error"
)

// RegisterRoutes register routes
func RegisterRoutes(e *gin.Engine) {
	e.POST("/admin/login", handleAcctLogin)
}

// RegisterRoutesAuthz register routes
func RegisterRoutesAuthz(group gin.IRoutes) {
	group.POST("/api/account", handleAPIAccount)
	group.POST("/api/blog", handleAPIBlogger)
	group.POST("/api/password", handleAPIPassword)
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

// handleAcctLogin 更新账户信息
func handleAPIAccount(c *gin.Context) {
	email := c.PostForm("email")
	phoneN := c.PostForm("phoneNumber")
	address := c.PostForm("address")
	if address == "" || !tools.ValidateEmail(email) || tools.ValidatePhoneNo(phoneN) {
		ResponseNotice(c, NoticeWarning, "参数异常", "")
		return
	}
	err := cache.Wi.Store.UpdateAccount(context.Background(), cache.Wi.Account.Username, map[string]interface{}{
		"email":   email,
		"address": address,
		"phone_n": phoneN,
	})
	if err != nil {
		logrus.Error("handleAPIAccount.UpdateAccount: ", err)
		ResponseNotice(c, NoticeWarning, err.Error(), "")
		return
	}
	cache.Wi.Account.Email = email
	cache.Wi.Account.Address = address
	cache.Wi.Account.PhoneN = phoneN
	ResponseNotice(c, NoticeSuccess, "更新成功", "")
}

// handleAPIBlogger 更新博客信息
func handleAPIBlogger(c *gin.Context) {
	blogName := c.PostForm("blogName")
	bTitle := c.PostForm("bTitle")
	subTitle := c.PostForm("subTitle")
	beiAn := c.PostForm("beiAn")
	seriessay := c.PostForm("seriessay")
	archivessay := c.PostForm("archivessay")
	if blogName == "" || bTitle == "" {
		ResponseNotice(c, NoticeWarning, "参数错误", "")
		return
	}
	err := cache.Wi.Store.UpdateBlogger(c, map[string]interface{}{
		"blog_name":    blogName,
		"b_title":      bTitle,
		"bei_an":       beiAn,
		"sub_title":    subTitle,
		"series_say":   seriessay,
		"archives_say": archivessay,
	})
	if err != nil {
		logrus.Error("handleAPIBlogger.UpdateBlogger: ", err)
		ResponseNotice(c, NoticeWarning, err.Error(), "")
		return
	}
	cache.Wi.Blogger.BlogName = blogName
	cache.Wi.Blogger.BTitle = bTitle
	cache.Wi.Blogger.BeiAn = beiAn
	cache.Wi.Blogger.SubTitle = subTitle
	cache.Wi.Blogger.SeriesSay = seriessay
	cache.Wi.Blogger.ArchivesSay = archivessay

	cache.PagesCh <- cache.PageSeries
	cache.PagesCh <- cache.PageArchive

	ResponseNotice(c, NoticeSuccess, "更新成功", "")
}

// handleAPIPassword 更新密码
func handleAPIPassword(c *gin.Context) {
	od := c.PostForm("old")
	nw := c.PostForm("new")
	cf := c.PostForm("confirm")
	if nw != cf {
		ResponseNotice(c, NoticeWarning, "两次输入的密码不正确", "")
		return
	}
	if !tools.ValidatePassword(nw) {
		ResponseNotice(c, NoticeWarning, "密码格式不正确", "")
		return
	}
	if cache.Wi.Account.Password != tools.EncryptPassword(cache.Wi.Account.Username, od) {
		ResponseNotice(c, NoticeWarning, "原始密码不正确", "")
		return
	}
	newPw := tools.EncryptPassword(cache.Wi.Account.Username, nw)
	err := cache.Wi.Store.UpdateAccount(c, cache.Wi.Account.Username, map[string]interface{}{
		"password": newPw,
	})
	if err != nil {
		logrus.Error("handleAPIPassword.UpdateAccount: ", err)
		ResponseNotice(c, NoticeWarning, err.Error(), "")
		return
	}
	cache.Wi.Account.Password = newPw
	ResponseNotice(c, NoticeSuccess, "更新成功", "")
}

// ResponseNotice response notice
func ResponseNotice(c *gin.Context, typ, content, hl string) {
	if hl != "" {
		c.SetCookie("notice_highlight", hl, 86400, "/", "", true, false)
	}
	c.SetCookie("notice_type", typ, 86400, "/", "", true, false)
	c.SetCookie("notice", fmt.Sprintf("[\"%s\"]", content), 86400, "/", "", true, false)
	c.Redirect(http.StatusFound, c.Request.Referer())
}
