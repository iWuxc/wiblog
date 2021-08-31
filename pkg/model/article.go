// Package model provides ...
package model

import (
	"time"
)

// use snake_case as column name

// Article 文章
type Article struct {
	ID      int    `gorm:"column:id;primaryKey" bson:"id"`                           // ID, store自行控制
	Author  string `gorm:"column:author;type:string;size:50;not null" bson:"author"` // 作者名
	Slug    string `gorm:"column:slug;type:string;size:250;not null" bson:"slug"`    // 文章缩略名
	Title   string `gorm:"column:title;type:string;size:250;not null" bson:"title"`  // 标题
	Count   int    `gorm:"column:count;not null" bson:"count"`                       // 评论数量
	Content string `gorm:"column:content;not null" bson:"content"`                   // markdown内容
	SerieID int    `gorm:"column:serie_id;not null" bson:"serie_id"`                 // 专题ID
	Tags    string `gorm:"column:tags;type:string;size:250" bson:"tags"`             // tags
	IsDraft bool   `gorm:"column:is_draft;not null" bson:"is_draft"`                 // 是否是草稿

	DeletedAt time.Time `gorm:"column:deleted_at;not null,index:index_deleted_at" bson:"deleted_at"` // 删除时间
	UpdatedAt time.Time `gorm:"column:updated_at" bson:"updated_at"`                                 // 更新时间
	CreatedAt time.Time `gorm:"column:created_at" bson:"created_at"`                                 // 创建时间

	Header  string   `gorm:"-" bson:"-"` // header
	Excerpt string   `gorm:"-" bson:"-"` // 预览信息
	Desc    string   `gorm:"-" bson:"-"` // 描述
	Thread  string   `gorm:"-" bson:"-"` // disqus thread
	Prev    *Article `gorm:"-" bson:"-"` // 上篇文章
	Next    *Article `gorm:"-" bson:"-"` // 下篇文章
}

// SortedArticles 按时间排序后文章
type SortedArticles []*Article

// Len 长度
func (s SortedArticles) Len() int { return len(s) }

// Less 对比
func (s SortedArticles) Less(i, j int) bool { return s[i].CreatedAt.After(s[j].CreatedAt) }

// Swap 交换
func (s SortedArticles) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
