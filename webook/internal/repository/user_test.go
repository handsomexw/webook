package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	cachemocks "basic-go/webook/internal/repository/cache/mocks"
	"basic-go/webook/internal/repository/dao"
	daomocks "basic-go/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	//去掉毫秒以外的部分
	now = time.UnixMilli(now.UnixMilli())
	testCase := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		ctx  context.Context
		id   int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，数据库查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExit)
				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "这是密码",
					Phone: sql.NullString{
						String: "13456789",
						Valid:  true,
					},
					Ctime: now.UnixMilli(),
					Utime: now.UnixMilli(),
				}, nil)
				uc.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "这是密码",
					Phone:    "13456789",
					Ctime:    now,
				}).Return(nil)
				return ud, uc
			},
			ctx: context.Background(),
			id:  int64(123),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "这是密码",
				Phone:    "13456789",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "这是密码",
					Phone:    "13456789",
					Ctime:    now,
				}, nil)
				ud := daomocks.NewMockUserDao(ctrl)
				return ud, uc
			},
			ctx: context.Background(),
			id:  int64(123),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "这是密码",
				Phone:    "13456789",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExit)
				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, dao.ErrorUserNotFound)
				return ud, uc
			},
			ctx:      context.Background(),
			id:       int64(123),
			wantUser: domain.User{},
			wantErr:  dao.ErrorUserNotFound,
		},
		{
			name: "缓存未命中，数据库查询成功，缓存刷新失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExit)
				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "这是密码",
					Phone: sql.NullString{
						String: "13456789",
						Valid:  true,
					},
					Ctime: now.UnixMilli(),
					Utime: now.UnixMilli(),
				}, nil)
				uc.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "这是密码",
					Phone:    "13456789",
					Ctime:    now,
				}).Return(errors.New("缓存错误"))
				return ud, uc
			},
			ctx: context.Background(),
			id:  int64(123),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "这是密码",
				Phone:    "13456789",
				Ctime:    now,
			},
			wantErr: errors.New("缓存错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			user, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}
}
