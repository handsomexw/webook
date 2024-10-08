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

type UserDao interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Insert(ctx context.Context, user User) error
	FindByOpenId(ctx context.Context, openId string) (User, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

type User struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`
	//唯一索引运行有多个空值
	Email         sql.NullString `gorm:"unique"`
	Password      string
	Phone         sql.NullString `gorm:"unique"`
	WechatUnionID sql.NullString
	WechatOpenID  sql.NullString `gorm:"unique"`
	Ctime         int64
	Utime         int64
}

type UserInfo struct {
	Name            string
	Birthday        string
	PersonalProfile string
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{db: db}
}

func (ud *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ud *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (ud *GORMUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("`Id` = ?", id).First(&user).Error
	return user, err
}

func (ud *GORMUserDao) Insert(ctx context.Context, user User) error {
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

func (ud *GORMUserDao) FindByOpenId(ctx context.Context, openId string) (User, error) {
	var user User
	//gorm默认列名规则
	err := ud.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&user).Error
	return user, err
}
