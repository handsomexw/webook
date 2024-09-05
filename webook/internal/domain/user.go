package domain

import "time"

type User struct {
	Id         int64
	Email      string
	Password   string
	Phone      string
	Ctime      time.Time
	Name       string
	AboutMe    string
	Birthday   time.Time
	WechatInfo WechatInfo
}
