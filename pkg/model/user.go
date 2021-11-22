// Package model provides ...
package model

import "time"

type User struct {
	ID          int       `gorm:"column:id;primaryKey" bson:"id" json:"id"`                            // 自增ID
	NickName    string    `gorm:"column:nickname" bson:"nick_name" json:"nickName"`                    // 昵称(QQ)
	Gender      string    `gorm:"column:gender;size:50" bson:"gender" json:"gender"`                   //性别
	Figure40    string    `gorm:"column:figureurl_40;size:255" bson:"figure_40" json:"figure_40"`       //40×40头像
	Figure100   string    `gorm:"column:figureurl_100;size:255" bson:"figure_100" json:"figure_100"`    //40×40头像
	OpenId      string    `gorm:"column:openid;size:255" bson:"openid" json:"openid"`                   //用户身份信息(QQ)
	AccessToken string    `gorm:"column:access_token;size:255" bson:"access_token" json:"access_token"` //accessToken(QQ)
	InvalidTime time.Time `gorm:"column:invalid_time" bson:"invalid_time" json:"invalid_time"`         //Token时效时间[前10分钟]

	DeletedAt time.Time `gorm:"column:deleted_at;not null,index:index_deleted_at" bson:"deleted_at" json:"deleted_at"` // 删除时间
	UpdatedAt time.Time `gorm:"column:updated_at" bson:"updated_at" json:"updated_at"`                                 // 更新时间
	CreatedAt time.Time `gorm:"column:created_at" bson:"created_at" json:"created_at"`                                 // 创建时间
}
