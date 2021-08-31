// Package store provides ...
package store

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	)

	db.DB = gormDB
	return db, nil
}

// LoadInsertAccount 读取或创建账户
func (db *rdbms) LoadInsertAccount(ctx context.Context, acct *model.Account) (bool, error) {
	result := db.Where("username=?", acct.Username).FirstOrCreate(acct)
	return result.RowsAffected > 0, result.Error
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
