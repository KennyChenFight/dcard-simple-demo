package model

import "time"

type User struct {
	Id             string    `xorm:"pk" json:"id" update:"fixed"`
	Email          string    `json:"email" binding:"required,email"`
	PasswordDigest string    `json:"-"`
	Name           string    `json:"name" binding:"required"`
	CreateTime     time.Time `json:"create_time" xorm:"created utc"`
	UpdateTime     time.Time `json:"update_time" xorm:"updated utc"`
}

func (u *User) TableName() string {
	return "users"
}
