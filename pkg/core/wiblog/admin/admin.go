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
	"wiblog/pkg/conf"
	"wiblog/pkg/core/wiblog"
	"wiblog/pkg/internal"
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
	group.POST("/api/post-delete", handleAPIPostDelete)
	group.POST("/api/serie-add", handleAPISerieCreate)
	group.POST("/api/serie-delete", handleAPISerieDelete)
	group.POST("/api/draft-delete", handleDraftDelete)
	group.POST("/api/trash-delete", handleAPITrashDelete)
	group.POST("/api/trash-recover", handleAPITrashRecover)
	group.POST("/api/file-upload", handleAPIQiniuUpload)
	group.POST("/api/file-delete", handleAPIQiniuDelete)
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

// handleAPIPostCreate 添加文章
func handleAPIPostCreate(c *gin.Context) {
	var (
		err error
		do  string
		cid int
	)

	defer func() {
		now := time.Now().In(tools.TimeLocation)
		switch do {
		case "auto": // 自动保存
			if nil != err {
				c.JSON(http.StatusOK, gin.H{"fail": 1, "time": now.Format("15:04:05 PM"), "cid": cid})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": 0, "time": now.Format("15:04:05 PM"), "cid": cid})
		case "save", "publish": // 草稿，发布
			if err != nil {
				ResponseNotice(c, NoticeWarning, err.Error(), "")
				return
			}
			uri := "/admin/manage-draft"
			if do == "publish" {
				uri = "/admin/manage-posts"
			}
			c.Redirect(http.StatusFound, uri)
		}
	}()

	do = c.PostForm("do")
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	text := c.PostForm("text")
	cover := c.PostForm("cover")
	isHot := c.PostForm("hot")

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
		Title:     title,
		Content:   text,
		Slug:      slug,
		Cover:     cover,
		IsDraft:   do != "publish",
		Author:    cache.Wi.Account.Username,
		SerieID:   serieid,
		Tags:      tags,
		CreatedAt: date,
		IsHot: isHot == "1",
	}

	cid, err = strconv.Atoi(c.PostForm("cid"))
	//新文章
	if err != nil || cid < 1 {
		err = cache.Wi.AddArticle(article)
		if err != nil {
			logrus.Error("handleAPIPostCreate.AddArticle: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}

		//TODO::异步执行其他文章操作
		return
	}
	//旧文章
	err = cache.Wi.Store.UpdateArticle(context.Background(), cid, map[string]interface{}{
		"title":      article.Title,
		"content":    article.Content,
		"cover":      article.Cover,
		"serie_id":   article.SerieID,
		"is_draft":   article.IsDraft,
		"tags":       article.Tags,
		"updated_at": article.UpdatedAt,
		"created_at": article.CreatedAt,
	})
	if err != nil {
		logrus.Error("handleAPIPostCreate.UpdateArticle: ", err)
		return
	}
	ResponseNotice(c, NoticeSuccess, "操作成功", "")
}

// handleAPIPostDelete 删除文章，移入回收箱
func handleAPIPostDelete(c *gin.Context) {
	var ids []int
	for _, v := range c.PostFormArray("cid[]") {
		fmt.Println(v)
		id, err := strconv.Atoi(v)
		if err != nil || id < conf.Conf.WiBlogApp.General.StartID {
			ResponseNotice(c, NoticeWarning, "参数错误", "")
			return
		}
		err = cache.Wi.DelArticle(id)
		if err != nil {
			logrus.Error("handleAPIPostDelete.DelArticle: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
		ids = append(ids, id)
	}
	ResponseNotice(c, NoticeSuccess, "删除成功", "")
	//TODO::elasticsearch
}

//handleDraftDelete 硬删除文章
func handleDraftDelete(c *gin.Context) {
	for _, v := range c.PostFormArray("mid[]") {
		id, err := strconv.Atoi(v)
		if err != nil || id < 0 {
			ResponseNotice(c, NoticeWarning, "参数错误", "")
			return
		}
		err = cache.Wi.Store.RemoveArticle(context.Background(), id)
		if err != nil {
			logrus.Error("handleDraftDelete.RemoveArticle: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
	}
	ResponseNotice(c, NoticeSuccess, "删除成功", "")
}

// handleAPITrashDelete 删除回收箱文章
func handleAPITrashDelete(c *gin.Context) {
	for _, v := range c.PostFormArray("mid[]") {
		id, err := strconv.Atoi(v)
		if err != nil || id < 0 {
			ResponseNotice(c, NoticeWarning, "参数错误", "")
			return
		}
		err = cache.Wi.Store.RemoveArticle(context.Background(), id)
		if err != nil {
			logrus.Error("handleAPITrashDelete.RemoveArticle: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
	}
	ResponseNotice(c, NoticeSuccess, "删除成功", "")
}

// handleAPITrashRecover 恢复到草稿箱
func handleAPITrashRecover(c *gin.Context) {
	for _, v := range c.PostFormArray("mid[]") {
		id, err := strconv.Atoi(v)
		if err != nil || id < 0 {
			ResponseNotice(c, NoticeWarning, "参数错误", "")
			return
		}
		err = cache.Wi.Store.UpdateArticle(context.Background(), id, map[string]interface{}{
			"deleted_at": time.Time{},
			"is_draft":   true,
		})
		if err != nil {
			logrus.Error("handleAPITrashRecover.UpdateArticle: ", err)
			ResponseNotice(c, NoticeWarning, err.Error(), "")
			return
		}
	}
	ResponseNotice(c, NoticeSuccess, "恢复成功", "")
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

// handleAPISerieDelete 专题删除
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

// handleAPIQiniuUpload 七牛云上传
func handleAPIQiniuUpload(c *gin.Context) {
	type Size interface {
		Size() int64
	}
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		logrus.Error("handleAPIQiniuUpload.FormFile: ", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	s, ok := file.(Size)
	if !ok {
		logrus.Error("assert failed ")
		c.String(http.StatusBadRequest, "false")
		return
	}
	filename := strings.ToLower(header.Filename)

	params := internal.UploadParams{
		Name: filename,
		Size: s.Size(),
		Data: file,

		Config: conf.Conf.WiBlogApp.Qiniu,
	}

	url, err := internal.QiniuUpload(params)
	if err != nil {
		logrus.Error("handleAPIQiniuUpload.QiniuUpload: ", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	typ := header.Header.Get("Content-Type")
	c.JSON(http.StatusOK, gin.H{
		"title":   filename,
		"isImage": typ[:5] == "image",
		"url":     url,
		"bytes":   fmt.Sprintf("%dkb", s.Size()/1000),
	})
}

// handleAPIQiniuDelete 七牛云删除
func handleAPIQiniuDelete(c *gin.Context) {
	defer c.String(http.StatusOK, "删掉了吗？鬼知道。。。")
	name := c.PostForm("title")
	if name == "" {
		logrus.Error("handleAPIQiniuDelete.PostForm: 参数错误")
		return
	}
	params := internal.DeleteParams{
		Name:   name,
		Config: conf.Conf.WiBlogApp.Qiniu,
	}
	err := internal.QiniuDelete(params)
	if err != nil {
		logrus.Error("handleAPIQiniuDelete.QiniuDelete: ", err)
	}
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
