package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrorUserDuplicate = errors.New("邮箱冲突")
	ErrorUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

type User struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`
	//唯一索引运行有多个空值
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
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

func (ud *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
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
		const uniqueConflictErrNo uint16 = 1062
		if mysqlError.Number == uniqueConflictErrNo {
			//邮箱或者手机号码冲突
			return ErrorUserDuplicate
		}
	}
	return err
}
