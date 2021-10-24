// Package model provides ...
package model

import (
	"github.com/lib/pq"
	"sync"
	"time"
)

// use snake_case as column name

// Article 文章
type Article struct {
	ID      int            `gorm:"column:id;primaryKey" bson:"id" json:"id"`                               // ID, store自行控制
	Author  string         `gorm:"column:author;type:string;size:50;not null" bson:"author" json:"author"` // 作者名
	Slug    string         `gorm:"column:slug;type:string;size:250;not null" bson:"slug" json:"slug"`      // 文章缩略名
	Title   string         `gorm:"column:title;type:string;size:250;not null" bson:"title" json:"title"`   // 标题
	Cover   string         `gorm:"column:cover;type:string;size:250;not null" bson:"cover" json:"cover"`   //封面图
	Count   int            `gorm:"column:count;not null" bson:"count" json:"count"`                        // 评论数量
	Content string         `gorm:"column:content;not null" bson:"content" json:"content"`                  // markdown内容
	SerieID int            `gorm:"column:serie_id;not null" bson:"serie_id" json:"serie_id"`               // 专题ID
	Tags    pq.StringArray `gorm:"column:tags;type:string;size:255;default:'{}'" bson:"tags" json:"tags"`  // tags
	IsDraft bool           `gorm:"column:is_draft;not null" bson:"is_draft" json:"is_draft"`               // 是否是草稿
	IsHot   bool           `gorm:"column:is_hot;not null" bson:"is_hot" json:"is_hot"`                     // 热门文章

	DeletedAt time.Time `gorm:"column:deleted_at;not null,index:index_deleted_at" bson:"deleted_at" json:"deleted_at"` // 删除时间
	UpdatedAt time.Time `gorm:"column:updated_at" bson:"updated_at" json:"updated_at"`                                 // 更新时间
	CreatedAt time.Time `gorm:"column:created_at" bson:"created_at" json:"created_at"`                                 // 创建时间

	Header  string   `gorm:"-" bson:"-" json:"-"` // header
	Excerpt string   `gorm:"-" bson:"-" json:"-"` // 预览信息
	Desc    string   `gorm:"-" bson:"-" json:"-"` // 描述
	Thread  string   `gorm:"-" bson:"-" json:"-"` // disqus thread
	Prev    *Article `gorm:"-" bson:"-" json:"-"` // 上篇文章
	Next    *Article `gorm:"-" bson:"-" json:"-"` // 下篇文章

	CreatedDay  string `gorm:"-" bson:"-" json:"created_day"`  //创建日期 - 天
	CreatedMon  string `gorm:"-" bson:"-" json:"created_mon"`  //创建日期 - 月
	CreatedYear string `gorm:"-" bson:"-" json:"created_year"` //创建日期 - 年
	SerieName   string `gorm:"-" bson:"-" json:"serie_name"`
	ArticleUrl  string `gorm:"-" bson:"-" json:"article_url"` //访问地址
}

// SortedArticles 按时间排序后文章
type SortedArticles []*Article

// ArticleList sync handle
type ArticleList struct {
	Lock  *sync.Mutex
	IdMap map[int]*Article
}

// Len 长度
func (s SortedArticles) Len() int { return len(s) }

// Less 对比
func (s SortedArticles) Less(i, j int) bool { return s[i].CreatedAt.After(s[j].CreatedAt) }

// Swap 交换
func (s SortedArticles) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
