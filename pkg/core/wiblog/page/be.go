// Package page provides ...
package page

import (
	"bytes"
	"github.com/gin-gonic/gin"
	htemplate "html/template"
	"net/http"
	"wiblog/pkg/cache"
	"wiblog/pkg/conf"
	"wiblog/pkg/core/wiblog"
)

// baseBEParams 页面基础参数
func baseBEParams(c *gin.Context) gin.H {
	return gin.H{
		"Author": cache.Wi.Account.Username,
		"Qiniu": conf.Conf.WiBlogApp.Qiniu,
	}
}

// handleLoginPage 登录页面
func handleLoginPage(c *gin.Context) {
	logout := c.Query("logout")
	if logout == "true" {
		wiblog.SetLogout(c)
	} else if wiblog.IsLogined(c) {
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}
	params := gin.H{"BTitle": cache.Wi.Blogger.BTitle}
	renderHTMLAdminLayout(c, "login.html", params)
}

// handleAdminProfile 个人配置中心
func handleAdminProfile(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "个人配置 | " + cache.Wi.Blogger.BTitle
	params["Path"] = c.Request.URL.Path
	params["Console"] = true
	params["Wi"] = cache.Wi
	renderHTMLAdminLayout(c, "admin-profile", params)
}

// renderHTMLAdminLayout 渲染admin页面
func renderHTMLAdminLayout(c *gin.Context, name string, data gin.H) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	// special page
	if name == "login.html" {
		err := htmlTemplate.ExecuteTemplate(c.Writer, name, data)
		if err != nil {
			panic(err)
		}
		return
	}
	buf := bytes.Buffer{}
	err := htmlTemplate.ExecuteTemplate(&buf, name, data)
	if err != nil {
		panic(err)
	}
	data["LayoutContent"] = htemplate.HTML(buf.String())
	err = htmlTemplate.ExecuteTemplate(c.Writer, "adminLayout.html", data)
	if err != nil {
		panic(err)
	}
	if c.Writer.Status() == 0 {
		c.Status(http.StatusOK)
	}
}
