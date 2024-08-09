package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var (
	ErrorUserDuplicateEmail = dao.ErrorUserDuplicateEmail
	ErrorUserNotFound       = dao.ErrorUserNotFound
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
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
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
	u := domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	//redis没找到，数据库找到后要放回redis
	err = ur.cache.Set(ctx, u)
	if err != nil {
		//打日志，做监控
	}
	return u, err

}
