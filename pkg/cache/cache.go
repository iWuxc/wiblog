// Package cache provides ...
package cache

import (
	"context"
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
	}

	//load and init
	err = Wi.loadOrInit()
	if err != nil {
		panic(err)
	}

}

type Cache struct {
	lock sync.Mutex
	Store store.Store

	//load from model
	Account *model.Account
	Blogger *model.Blogger
	Articles *model.Article
}

func (c *Cache) loadOrInit() error {
	blogapp := conf.Conf.WiBlogApp

	blogger := &model.Blogger{
		BlogName: strings.Title(blogapp.Account.Username),
		SubTitle: "Rome was not built in one day.",
		BeiAn: "蜀ICP备xxxxxxxx号-1",
		BTitle: fmt.Sprintf("%s's Blog", strings.Title(blogapp.Account.Username)),
		Copyright: `本站使用「<a href="//creativecommons.org/licenses/by/4.0/">署名 4.0 国际</a>」创作共享协议，转载请注明作者及原网址。`,
	}

	created, err := c.Store.LoadInsertBlogger(context.Background(), blogger)
	if err != nil {
		return err
	}

	c.Blogger = blogger
	if created { // init articles: about blogroll
		about := &model.Article{
			ID: 1, //固定ID
			Author: blogapp.Account.Username,
			Title: "关于",
			Slug: "about",
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

	return nil
}
