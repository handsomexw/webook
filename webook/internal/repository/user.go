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

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (ur *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDoamin(u), nil
}

func (ur *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDoamin(u), nil
}

func (ur *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToEntity(u))
}

func (ur *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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

func (ur *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (ur *CacheUserRepository) entityToDoamin(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		WechatInfo: domain.WechatInfo{
			OpenID:  u.WechatOpenID.String,
			UnionID: u.WechatUnionID.String,
		},
		Ctime: time.UnixMilli(u.Ctime),
	}
}

func (ur *CacheUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := ur.dao.FindByOpenId(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDoamin(u), nil
}
