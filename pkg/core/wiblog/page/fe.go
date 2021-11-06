package page

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	htemplate "html/template"
	"net/http"
	"strings"
	"time"
	"wiblog/pkg/cache"
	"wiblog/pkg/cache/render"
	"wiblog/pkg/cache/store"
	"wiblog/pkg/conf"
	"wiblog/pkg/service"
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
	search := store.SearchArticles{
		Page: 1,
		Limit: 3,
		Fields: map[string]interface{}{
			store.SearchArticleHot: true,
		},
	}
	articles, _, _ := service.ArticleList(c, search)
	params["Article"] = articles
	renderHTMLHomeLayout(c, "home.html", params)
}

// handleArticleIndexPage 文章列表
func handleArticleIndexPage(c *gin.Context) {

	params := baseFEParams(c)
	params["Title"] = "文章列表" + " | " + cache.Wi.Blogger.SubTitle
	params["Description"] = "文章列表，" + cache.Wi.Blogger.SubTitle
	params["CurrentPage"] = "web-posts"

	//文章设置属性
	pagesize := conf.Conf.WiBlogApp.General.PageNum
	params["PageSize"] = pagesize

	var search = store.SearchArticles{
		Page:  1,
		Limit: 10,
		Fields: map[string]interface{}{
			store.SearchArticleDraft: false,
		},
	}

	count, _ := cache.Wi.Store.LoadArticleCount(c, search)
	params["Count"] = count
	//总页数
	params["PageCount"] = (count + pagesize - 1) / pagesize

	//分类
	params["Series"] = cache.Wi.Series

	renderHTMLHomeLayout(c, "web-posts", params)
}

// handleArticleDetailPage 文章详情
func handleArticleDetailPage(c *gin.Context) {
	params := baseFEParams(c)

	params["Title"] = "文章详情" + " | " + cache.Wi.Blogger.SubTitle
	params["Description"] = "文章详情，" + cache.Wi.Blogger.SubTitle
	params["CurrentPage"] = "web-posts"

	slug := c.Param("slug")

	if !strings.HasSuffix(slug, ".html") {
		renderHTMLHomeLayout(c, "404.html", params)
		return
	}

	article, err := cache.Wi.Store.FindArticleBySlug(c, slug[:len(slug)-5])
	if err != nil || article == nil {
		renderHTMLHomeLayout(c, "404.html", params)
		return
	}

	render.GenerateExcerptMarkdown(article)

	params["Article"] = article
	renderHTMLHomeLayout(c, "web-post", params)
}

// renderHTMLHomeLayout homelayout html
func renderHTMLHomeLayout(c *gin.Context, name string, data gin.H) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	//special page
	if name == "home.html" || name == "404.html" {
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
