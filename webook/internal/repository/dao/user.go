package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrorUserDuplicateEmail = errors.New("邮箱冲突")
	ErrorUserNotFound       = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}

type UserInfo struct {
	Name            string
	Birthday        string
	PersonalProfile string
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (ud *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ud *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("`Id` = ?", id).First(&user).Error
	return user, err
}

func (ud *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now

	err := ud.db.WithContext(ctx).Create(&user).Error
	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) {
		const uniqueConflictEmail uint16 = 1062
		if mysqlError.Number == uniqueConflictEmail {
			return ErrorUserDuplicateEmail
		}
	}
	return err
}
