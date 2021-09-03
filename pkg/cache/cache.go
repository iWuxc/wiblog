// Package cache provides ...
package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
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
		lock:  sync.Mutex{},
		Store: s,
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
	Articles *model.Article

	//auto genernal
	PageSeries  string //page
	Series      model.SortedSeries
	TagArticles map[string]model.SortedArticles // tagname:articles
	ArticlesMap map[string]*model.Article       // slug:article
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

	//series init
	series, err := c.Store.LoadAllSerie(context.Background())
	if err != nil {
		return nil
	}
	c.Series = series

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
