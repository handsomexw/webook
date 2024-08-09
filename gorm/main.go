package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Name     string `gorm:"primaryKey"`
	Age      int
	Birthday string
}

func main() {
	//dsn := "root:123456@tcp(127.0.0.1:13306)/test?charset=utf8&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	DryRun: true,
	//})
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	type Result struct {
		Name string
		Age  int
	}
	re := Result{}
	us := User{}
	var users []User
	err = db.AutoMigrate(&User{})
	if err != nil {
		return
	}
	db.Create(&us)
	db.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "a%").Group("name").Find(&re)
	db.Model(&us).Where("age > ?", "10").Updates(map[string]interface{}{"age": 10})
	db.Model(&User{}).Select("user.name, emails.email").Joins("left join emails on emails.uer_id = users.is").Scan(&re)
	db.Where("name = ? AND age > ?", "a%", 10).Find(&users)
	db.Find(&users, "age > ?", 10)
	db.Delete(&User{}, "name = ? AND age > ?", "a%", 10)
	time.Sleep(time.Second)
}
