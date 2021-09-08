package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wiblog/pkg/cache"
	"wiblog/pkg/core/wiblog"
	"wiblog/pkg/model"
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
	group.POST("/api/post-add", handleAPIPostCreate)
	group.POST("/api/serie-add", handleAPISerieCreate)
	group.POST("/api/serie-delete", handleAPISerieDelete)
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
	//
	//cache.PagesCh <- cache.PageSeries
	//cache.PagesCh <- cache.PageArchive

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

func handleAPIPostCreate(c *gin.Context) {

	do := c.PostForm("do")
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	text := c.PostForm("text")
	//cid := c.PostForm("cid")
	date := parseLocationDate(c.PostForm("title"))
	serie := c.PostForm("serie")
	tag := c.PostForm("tags")
	if title == "" || slug == "" || text == "" {
		ResponseNotice(c, NoticeWarning, "参数错误", "")
		return
	}
	var tags []string
	if tag != "" {
		tags = strings.Split(tag, ",")
	}
	serieid, _ := strconv.Atoi(serie)
	article := &model.Article{
		Title: title,
		Content: text,
		Slug: slug,
		IsDraft: do != "publish",
		Author: cache.Wi.Account.Username,
		SerieID: serieid,
		Tags:      tags,
		CreatedAt: date,
	}

	fmt.Println(article)
	err := cache.Wi.Store.InsertArticle(context.Background(), article, cache.ArticleStartID)
	if err != nil {
		logrus.Error("handleAPIPostCreate.InsertArticle: ", err)
		ResponseNotice(c, NoticeWarning, err.Error(), "")
		return
	}
	ResponseNotice(c, NoticeSuccess, "操作成功", "")
}

// handleAPISerieCreate 创建专题 如果专题有提交 mid 即更新专题
func handleAPISerieCreate(c *gin.Context) {
	name := c.PostForm("name")
	slug := c.PostForm("slug")
	desc := c.PostForm("description")
	if name == "" || slug == "" || desc == "" {
		ResponseNotice(c, NoticeWarning, "参数错误", "")
		return
	}
	sid, err := strconv.Atoi(c.PostForm("sid"))
	if err == nil && sid > 0 {
		var serie *model.Serie
		for _, v := range cache.Wi.Series {
			if v.ID == sid {
				serie = v
				break
			}
		}
		if serie == nil {
			ResponseNotice(c, NoticeWarning, "专题不存在", "")
			return
		}
		err = cache.Wi.Store.UpdateSerie(c, sid, map[string]interface{}{
			"name": name,
			"slug": slug,
			"desc": desc,
		})
		if err != nil {
			logrus.Error("handleAPISerieCreate.UpdateSerie: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
		serie.Name = name
		serie.Slug = slug
		serie.Desc = desc
	} else {
		err = cache.Wi.AddSerie(&model.Serie{
			Name: name,
			Slug: slug,
			Desc: desc,
		})
		if err != nil {
			logrus.Error("handleAPISerieCreate.InsertSerie: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
	}
	ResponseNotice(c, NoticeSuccess, "操作成功", "")
}

// handleAPISerieDelete
func handleAPISerieDelete(c *gin.Context) {
	fmt.Println(c.PostFormArray("sid[]"))
	for _, v := range c.PostFormArray("sid[]") {
		id, err := strconv.Atoi(v)
		if err != nil && id < 1 {
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
		err = cache.Wi.DelSerie(id)
		if err != nil {
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
	}
	ResponseNotice(c, NoticeSuccess, "删除成功", "")
}

// parseLocationDate 解析日期
func parseLocationDate(date string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04", date, tools.TimeLocation)
	if err == nil {
		return t.UTC()
	}
	return time.Now()
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
