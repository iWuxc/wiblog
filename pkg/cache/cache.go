// Package cache provides ...
package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
	"sync"
	"time"
	"wiblog/pkg/cache/render"
	"wiblog/pkg/cache/store"
	"wiblog/pkg/conf"
	"wiblog/pkg/model"
	"wiblog/tools"
)

var (
	// Wi wiblog cache
	Wi *Cache

	// regenerate pages chan
	PagesCh     = make(chan string, 2)
	PageSeries  = "series-md"
	PageArchive = "archive-md"

	// ArticleStartID article start id
	ArticleStartID = 11
)

// Init 数据库初始化, 建表, 加索引操作等
// name 应该为具体的关系数据库驱动名
func init() {
	var err error
	tools.TimeLocation, err = time.LoadLocation(conf.Conf.WiBlogApp.General.Timezone)
	if err != nil {
		panic(err)
	}

	//init store
	logrus.Info("store divers: ", store.Drivers())
	s, err := store.NewStore(conf.Conf.Database.Driver, conf.Conf.Database.Source)
	if err != nil {
		panic(err)
	}

	//wi init
	Wi = &Cache{
		lock:        sync.Mutex{},
		Store:       s,
		TagArticles: make(map[string]model.SortedArticles),
		ArticlesMap: make(map[string]*model.Article),
	}

	//load and init
	err = Wi.loadOrInit()
	if err != nil {
		panic(err)
	}

}

type Cache struct {
	lock  sync.Mutex
	Store store.Store

	//load from model
	Account  *model.Account
	Blogger  *model.Blogger
	Articles model.SortedArticles

	//auto genernal
	PageSeries  string //page
	Series      model.SortedSeries
	TagArticles map[string]model.SortedArticles // tagname:articles
	ArticlesMap map[string]*model.Article       // slug:article
}

func (c *Cache) readdArticle(article *model.Article, needSort bool) {
	//tag
	for _, tag := range article.Tags {
		if needSort {
			sort.Sort(c.TagArticles[tag])
		}
	}

	//series
	for i, serie := range c.Series {
		if serie.ID != article.SerieID {
			continue
		}
		c.Series[i].Articles = append(c.Series[i].Articles, article)
		if needSort {
			sort.Sort(c.Series[i].Articles)
			PagesCh <- PageSeries //重建专题
		}
	}
}

// loadOrInit Init 数据库初始化, 建表, 加索引操作等
func (c *Cache) loadOrInit() error {
	blogapp := conf.Conf.WiBlogApp

	blogger := &model.Blogger{
		BlogName:  strings.Title(blogapp.Account.Username),
		SubTitle:  "Rome was not built in one day.",
		BeiAn:     "蜀ICP备xxxxxxxx号-1",
		BTitle:    fmt.Sprintf("%s's Blog", strings.Title(blogapp.Account.Username)),
		Copyright: `本站使用「<a href="//creativecommons.org/licenses/by/4.0/">署名 4.0 国际</a>」创作共享协议，转载请注明作者及原网址。`,
	}

	created, err := c.Store.LoadInsertBlogger(context.Background(), blogger)
	if err != nil {
		return err
	}

	c.Blogger = blogger
	if created { // init articles: about blogroll
		about := &model.Article{
			ID:        1, //固定ID
			Author:    blogapp.Account.Username,
			Title:     "关于",
			Slug:      "about",
			CreatedAt: time.Time{}.AddDate(0, 0, 1),
		}
		err = c.Store.InsertArticle(context.Background(), about, 0)
		if err != nil {
			return err
		}

		blogroll := &model.Article{
			ID:        2, // 固定ID
			Author:    blogapp.Account.Username,
			Title:     "友情链接",
			Slug:      "blogroll",
			CreatedAt: time.Time{}.AddDate(0, 0, 1),
		}
		err = c.Store.InsertArticle(context.Background(), blogroll, ArticleStartID)
		if err != nil {
			return err
		}
	}

	//account init
	pwd := tools.EncryptPassword(blogapp.Account.Username, blogapp.Account.Password)
	account := &model.Account{
		Username: blogapp.Account.Username,
		Password: pwd,
	}
	_, err = c.Store.LoadInsertAccount(context.Background(), account)
	if err != nil {
		return err
	}
	//load account
	c.Account = account

	//series init
	series, err := c.Store.LoadAllSerie(context.Background())
	if err != nil {
		return err
	}
	c.Series = series

	//all articles
	search := store.SearchArticles{
		Page:   1, //当前页码
		Limit:  9999,                                                    //每页大小
		Fields: map[string]interface{}{store.SearchArticleDraft: false}} //字段:值

	articles, _, err := c.Store.LoadArticleList(context.Background(), search)

	for i, article := range articles {
		// 渲染页面
		render.GenerateExcerptMarkdown(article)
		c.ArticlesMap[article.Slug] = article
		// 分析文章
		if article.ID < ArticleStartID {
			continue
		}
		if i > 0 {
			article.Prev = articles[i-1]
		}
		if i < len(articles)-1 &&
			articles[i+1].ID >= ArticleStartID {
			article.Next = articles[i+1]
		}
	}

	Wi.Articles = articles
	// 重建专题与归档
	PagesCh <- PageSeries
	PagesCh <- PageArchive
	return nil
}

// AddSerie 添加专题
func (c *Cache) AddSerie(serie *model.Serie) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.Store.InsertSerie(context.Background(), serie)
	if err != nil {
		return err
	}
	c.Series = append(c.Series, serie)
	PagesCh <- PageSeries
	return nil
}

// DelSerie 删除专题
func (c *Cache) DelSerie(id int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i, serie := range c.Series {
		if serie.ID == id {
			if len(serie.Articles) > 0 {
				return errors.New("请删除改专题下的所有文章")
			}
			err := c.Store.RemoveSerie(context.Background(), id)
			if err != nil {
				return nil
			}
			c.Series[i] = nil
			c.Series = append(c.Series[:i], c.Series[i+1:]...)
			PagesCh <- PageSeries
			break
		}
	}
	return nil
}
