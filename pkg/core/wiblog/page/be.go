// Package page provides ...
package page

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	htemplate "html/template"
	"net/http"
	"strconv"
	"wiblog/pkg/cache"
	"wiblog/pkg/cache/store"
	"wiblog/pkg/conf"
	"wiblog/pkg/core/wiblog"
)

// baseBEParams 页面基础参数
func baseBEParams(c *gin.Context) gin.H {
	return gin.H{
		"Author": cache.Wi.Account.Username,
		"Qiniu":  conf.Conf.WiBlogApp.Qiniu,
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

type T struct {
	ID   string `json:"id"`
	Tags string `json:"tags"`
}

// handleAdminWritePost 编辑文章
func handleAdminWritePost(c *gin.Context) {
	params := baseBEParams(c)
	id, err := strconv.Atoi(c.Query("cid"))
	if err == nil && id > 0 {
		article, _ := cache.Wi.Store.LoadArticle(c, id)
		if article != nil {
			params["Title"] = "编辑文章"
			params["Edit"] = article
		}
	}
	if params["Title"] == nil {
		params["Title"] = "创建文章 | " + cache.Wi.Blogger.BTitle
	}
	params["Path"] = c.Request.URL.Path
	params["Domain"] = conf.Conf.WiBlogApp.Host
	params["Series"] = cache.Wi.Series
	var tags []T
	for tag := range cache.Wi.TagArticles {
		tags = append(tags, T{tag, tag})
	}
	str, _ := json.Marshal(tags)
	params["Tags"] = string(str)
	renderHTMLAdminLayout(c, "admin-post", params)
}

// handleAdminPosts 文章列表
func handleAdminPosts(c *gin.Context) {
	//receive search value
	serie := c.Query("serie")
	keywords := c.Query("keywords")
	serieId, err := strconv.Atoi(serie)
	if err != nil || serieId < 1 {
		serieId = 0
	}
	pg, err := strconv.Atoi(c.Query("page"))
	if err != nil || pg < 1 {
		pg = 1
	}

	//all request url query
	vals := c.Request.URL.Query()
	fmt.Println(vals)

	params := baseBEParams(c)
	params["Title"] = "文章管理 | " + cache.Wi.Blogger.BTitle
	params["Manage"] = true
	params["Path"] = c.Request.URL.Path
	params["Domain"] = conf.Conf.WiBlogApp.Host
	params["Series"] = cache.Wi.Series

	params["Serie"] = serieId
	params["KW"] = keywords

	var max int
	params["List"], max = cache.Wi.PageArticleBE(serieId, keywords, false, false, pg,
		conf.Conf.WiBlogApp.General.PageNum)
	if pg < max {
		vals.Set("page", fmt.Sprint(pg+1))
		params["Next"] = vals.Encode()
	}
	if pg > 1 {
		vals.Set("page", fmt.Sprint(pg-1))
		params["Prev"] = vals.Encode()
	}
	params["PP"] = make(map[int]string, max)
	for i := 0; i < max; i++ {
		vals.Set("page", fmt.Sprint(i+1))
		params["PP"].(map[int]string)[i+1] = vals.Encode()
	}
	params["Cur"] = pg

	renderHTMLAdminLayout(c, "admin-posts", params)
}

// handleAdminSeries 专题列表
func handleAdminSeries(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "专题管理 | " + cache.Wi.Blogger.BTitle
	params["Manage"] = true
	params["Path"] = c.Request.URL.Path
	params["List"] = cache.Wi.Series
	renderHTMLAdminLayout(c, "admin-series", params)
}

// handleAdminSerie 专题添加
func handleAdminSerie(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "新增专题 | " + cache.Wi.Blogger.BTitle
	id, err := strconv.Atoi(c.Query("sid"))
	if err == nil && id > 0 {
		for _, v := range cache.Wi.Series {
			if v.ID == id {
				params["Title"] = "编辑专题 | " + cache.Wi.Blogger.BTitle
				params["Edit"] = v
			}
		}
	}
	params["Manage"] = true
	params["Path"] = c.Request.URL.Path
	renderHTMLAdminLayout(c, "admin-serie", params)
}

// handleAdminTags 标签管理
func handleAdminTags(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "标签管理 | " + cache.Wi.Blogger.BTitle
	params["Manage"] = true
	params["Path"] = c.Request.URL.Path
	params["List"] = cache.Wi.TagArticles
	renderHTMLAdminLayout(c, "admin-tags", params)
}

// handleAdminDraft 草稿箱
func handleAdminDraft(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "草稿箱 | " + cache.Wi.Blogger.BTitle
	params["Manage"] = true
	params["Path"] = c.Request.URL.Path
	search := store.SearchArticles{
		Page:  1,
		Limit: 9999,
		Fields: map[string]interface{}{
			store.SearchArticleDraft: true,
		},
	}
	var err error
	params["List"], _, err = cache.Wi.Store.LoadArticleList(context.Background(), search)
	if err != nil {
		logrus.Error("handleAdminDraft.LoadArticleList: ", err)
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
	renderHTMLAdminLayout(c, "admin-draft", params)
}

// handleAdminTrash 回收箱
func handleAdminTrash(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "回收箱 | " + cache.Wi.Blogger.BTitle
	params["Manage"] = true
	params["Path"] = c.Request.URL.Path
	search := store.SearchArticles{
		Page:  1,
		Limit: 9999,
		Fields: map[string]interface{}{
			store.SearchArticleTrash: true,
		},
	}
	var err error
	params["List"], _, err = cache.Wi.Store.LoadArticleList(context.Background(), search)
	if err != nil {
		logrus.Error("handleAdminDraft.LoadArticleList: ", err)
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
	renderHTMLAdminLayout(c, "admin-trash", params)
}

// handleAdminGeneral 基本设置
func handleAdminGeneral(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "基本设置 | " + cache.Wi.Blogger.BTitle
	params["Setting"] = true
	params["Path"] = c.Request.URL.Path
	renderHTMLAdminLayout(c, "admin-general", params)
}

// handleAdminDiscussion 阅读设置
func handleAdminDiscussion(c *gin.Context) {
	params := baseBEParams(c)
	params["Title"] = "阅读设置 | " + cache.Wi.Blogger.BTitle
	params["Setting"] = true
	params["Path"] = c.Request.URL.Path
	renderHTMLAdminLayout(c, "admin-discussion", params)
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
