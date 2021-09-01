// Package store provides ...
package store

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"wiblog/pkg/model"
)

var (
	storeMu sync.RWMutex
	stores  = make(map[string]Driver)
)

type Store interface {

	// LoadInsertBlogger 读取或创建博客
	LoadInsertBlogger(ctx context.Context, blogger *model.Blogger) (bool, error)
	// UpdateBlogger 更新博客
	UpdateBlogger(ctx context.Context, fields map[string]interface{}) error

	// LoadInsertAccount 读取或创建用户
	LoadInsertAccount(ctx context.Context, acct *model.Account) (bool, error)
	// UpdateAccount 更新账户
	UpdateAccount(ctx context.Context, name string, fields map[string]interface{}) error

	// InsertArticle 创建文章
	InsertArticle(ctx context.Context, article *model.Article, startID int) error
}

type Driver interface {
	Init(name, source string) (Store, error)
}

// Register 注册驱动
func Register(name string, driver Driver) {
	storeMu.Lock()
	defer storeMu.Unlock()
	if driver == nil {
		panic("store: register driver is nil")
	}
	if _, ok := stores[name]; ok {
		panic("store: register called twice for driver " + name)
	}
	stores[name] = driver
}

//Drivers 获取所有
func Drivers() []string {
	storeMu.Lock()
	defer storeMu.Unlock()
	list := make([]string, len(stores))
	for name := range stores {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// NewStore 新建存储
func NewStore(name string, source string) (Store, error) {
	storeMu.RLock()
	driver, ok := stores[name]
	storeMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown driver %q (forgotten import?)", name)
	}
	return driver.Init(name, source)
}
