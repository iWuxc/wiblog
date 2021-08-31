// Package page provides ...
package page

import (
	"bytes"
	"github.com/gin-gonic/gin"
	htemplate "html/template"
	"net/http"
	"wiblog/pkg/cache"
	"wiblog/pkg/core/wiblog"
)


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
	//params := gin.H{"BTitle": "武晓晨"}
	renderHTMLAdminLayout(c, "login.html", params)
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
