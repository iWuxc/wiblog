// Package store provides ...
package store

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"wiblog/pkg/model"
)

type rdbms struct {
	*gorm.DB
}

//数据库驱动初始化
func init() {
	Register("mysql", &rdbms{})
	//TODO:: 注册其他驱动
}

func (db *rdbms) Init(name, source string) (Store, error) {
	var (
		gormDB *gorm.DB
		err    error
	)
	switch name {
	case "mysql":
		gormDB, err = gorm.Open(mysql.Open(source), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}

	_ = gormDB.AutoMigrate(
		&model.Account{},
		&model.Blogger{},
		&model.Article{},
		&model.Serie{},
	)

	db.DB = gormDB
	return db, nil
}

// LoadInsertAccount 读取或创建账户
func (db *rdbms) LoadInsertAccount(ctx context.Context, acct *model.Account) (bool, error) {
	result := db.Where("username=?", acct.Username).FirstOrCreate(acct)
	return result.RowsAffected > 0, result.Error
}

// UpdateBlogger 更新博客
func (db *rdbms) UpdateBlogger(ctx context.Context, fields map[string]interface{}) error {
	return db.Model(&model.Blogger{}).Session(&gorm.Session{AllowGlobalUpdate: true}).Updates(fields).Error
}

// UpdateAccount 更新账户
func (db *rdbms) UpdateAccount(ctx context.Context, name string, fields map[string]interface{}) error {
	return db.Model(&model.Account{}).Where("username=?", name).Updates(fields).Error
}

// LoadInsertBlogger 读取获创建博客
func (db *rdbms) LoadInsertBlogger(ctx context.Context, blogger *model.Blogger) (bool, error) {
	result := db.FirstOrCreate(blogger)
	return result.RowsAffected > 0, result.Error
}

// LoadAllSerie 读取所有的专题
func (db *rdbms) LoadAllSerie(ctx context.Context) (model.SortedSeries, error) {
	var series model.SortedSeries
	err := db.Order("id DESC").Find(&series).Error
	return series, err
}

// InsertArticle 创建文章
func (db *rdbms) InsertArticle(ctx context.Context, article *model.Article, startID int) error {
	if article.ID == 0 {
		// auto generate id
		var id int
		err := db.Model(model.Article{}).Select("Max(id)").Row().Scan(&id)
		if err != nil {
			return err
		}
		if id < startID {
			id = startID
		} else {
			id += 1
		}
		article.ID = id
	}
	return db.Create(article).Error
}

// LoadArticle 加载文章
func (db *rdbms) LoadArticle(ctx context.Context, id int) (*model.Article, error) {
	article := &model.Article{}
	err := db.Model(&model.Article{}).Where("id=?", id).First(article).Error
	return article, err
}

// LoadArticleList 查找文章列表
func (db *rdbms) LoadArticleList(ctx context.Context, search SearchArticles) (model.SortedArticles, int, error) {
	gormDB := db.Model(model.Article{})
	for k, v := range search.Fields {
		switch k {
		case SearchArticleDraft:
			if ok := v.(bool); ok {
				gormDB = gormDB.Where("is_draft=?", true)
			} else {
				gormDB = gormDB.Where("is_draft=? AND deleted_at=?", false, time.Time{})
			}
		case SearchArticleTrash:
			gormDB = gormDB.Where("deleted_at!=?", time.Time{})
		case SearchArticleTitle:
			gormDB = gormDB.Where("title LIKE ?", "%"+v.(string)+"%")
		case SearchArticleSerieID:
			gormDB = gormDB.Where("serie_id=?", v.(int))
		}
	}
	// search count
	var count int64
	err := gormDB.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var articles model.SortedArticles
	err = gormDB.Limit(search.Limit).Offset((search.Page - 1) * search.Page).Order("created_at DESC").
		Find(articles).Error
	if err != nil {
		return nil, 0, err
	}
	return articles, int(count), nil
}
