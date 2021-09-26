package page

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	htemplate "html/template"
	"net/http"
	"time"
	"wiblog/pkg/cache"
	"wiblog/pkg/conf"
)

// baseFEParams 前端基础参数
func baseFEParams(c *gin.Context) gin.H {
	version := 0
	cookie, err := c.Request.Cookie("v")
	if err != nil || cookie.Value != fmt.Sprint(conf.Conf.WiBlogApp.StaticVersion) {
		version = conf.Conf.WiBlogApp.StaticVersion
	}

	keyword := cache.Wi.Blogger.BlogName

	return gin.H{
		"BlogName":  cache.Wi.Blogger.BlogName,
		"SubTitle":  cache.Wi.Blogger.SubTitle,
		"BTitle":    cache.Wi.Blogger.BTitle,
		"BeiAn":     cache.Wi.Blogger.BeiAn,
		"Copyright": cache.Wi.Blogger.Copyright,
		"Domain":    conf.Conf.WiBlogApp.Host,
		"CopyYear":  time.Now().Year(),
		"Qiniu":     conf.Conf.WiBlogApp.Qiniu,
		"Version":   version,
		"KeyWord":   keyword,
	}
}

// handleNotFound not found page
func handleNotFound(c *gin.Context) {
	//params := baseFEParams(c)
	//params["Title"] = "Not Found"
	//params["Description"] = "404 Not Found"
	//params["Path"] = ""
	//c.Status(http.StatusNotFound)
	//renderHTMLHomeLayout(c, "notfound", params)
}

// handleHomePage 博客首页
func handleHomePage(c *gin.Context) {
	params := baseFEParams(c)
	params["Title"] = "博客首页" + " | " + cache.Wi.Blogger.SubTitle
	params["Description"] = "博客首页，" + cache.Wi.Blogger.SubTitle
	params["Address"] = "北京市海淀区"
	params["QQ"] = "1272105573"
	params["Email"] = "wuxc.ent@gmail.com"
	renderHTMLHomeLayout(c, "home.html", params)
}

// renderHTMLHomeLayout homelayout html
func renderHTMLHomeLayout(c *gin.Context, name string, data gin.H) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	fmt.Println(name)

	//special page
	if name == "home.html" {
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
	err = htmlTemplate.ExecuteTemplate(c.Writer, "homeLayout.html", data)
	if err != nil {
		panic(err)
	}
	if c.Writer.Status() == 0 {
		c.Status(http.StatusOK)
	}
}
