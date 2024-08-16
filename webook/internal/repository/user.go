package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

var (
	ErrorUserDuplicate = dao.ErrorUserDuplicate
	ErrorUserNotFound  = dao.ErrorUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDoamin(u), nil
}

func (ur *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDoamin(u), nil
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToEntity(u))
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	user, err := ur.cache.Get(ctx, id)
	if err == nil {
		return user, nil
	}

	//if errors.Is(err, cache.ErrKeyNotExit) {
	//	//去数据库里查找
	//}
	//redis出错，是否加载
	//选加载：保护数据库，限流
	//选不加载，用户体验差一点
	ue, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u := ur.entityToDoamin(ue)
	//redis没找到，数据库找到后要放回redis
	err = ur.cache.Set(ctx, u)
	if err != nil {
		//打日志，做监控
	}
	return u, err

}

func (ur *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  true,
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (ur *UserRepository) entityToDoamin(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
