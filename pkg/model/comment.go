// Package model provides ...
package model

import "time"

type Comment struct {
	ID        int    `gorm:"column:id;primaryKey" bson:"id" json:"id"`                         // ID, store自行控制
	UserId    int    `gorm:"column:user_id; default:0" bson:"user_id" json:"user_id"`          //用户ID
	ReplayId  int    `gorm:"column:replay_id; default:0" bson:"replay_id" json:"replay_id"`    //回复评论的ID
	ParentId  int    `gorm:"column:parent_id; default:0" bson:"parent_id" json:"parent_id"`    //父及评论ID
	ArticleId int    `gorm:"column:article_id; not null" bson:"article_id" json:"article_id"`  //文章ID
	Content   string `gorm:"column:content; not null; size:300" bson:"content" json:"content"` //文章内容
	Status    bool   `gorm:"column:status; default: 0;" bson:"status" json:"status"`           //审核状态
	IsLatest  bool   `gorm:"column:is_latest; default: 1;" bson:"is_latest" json:"is_latest"`  //是否最近留言

	DeletedAt time.Time `gorm:"column:deleted_at;not null,index:index_deleted_at" bson:"deleted_at" json:"deleted_at"` // 删除时间
	UpdatedAt time.Time `gorm:"column:updated_at" bson:"updated_at" json:"updated_at"`                                 // 更新时间
	CreatedAt time.Time `gorm:"column:created_at" bson:"created_at" json:"created_at"`                                 // 创建时间

}
