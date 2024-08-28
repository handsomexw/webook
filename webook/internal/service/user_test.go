package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	repomocks "basic-go/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

var now = time.Now()

func TestOneUserService_Login(t *testing.T) {
	testCase := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		ctx  context.Context
		user domain.User

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						//密码应该是加密的
						Password: "$2a$10$AVk680PliN79.CF3C.KFQeric11w2scHc0slHrK9B9WjPxRTHJta6",
						Phone:    "123456",
						Ctime:    now,
					}, nil)
				return repo
			},
			ctx: context.Background(),
			user: domain.User{
				Email: "123@qq.com",
				//密码应该是加密的
				Password: "123",
				Phone:    "123456",
				Ctime:    time.Now(),
			},
			wantUser: domain.User{
				Email: "123@qq.com",
				//密码应该是加密的
				Password: "$2a$10$AVk680PliN79.CF3C.KFQeric11w2scHc0slHrK9B9WjPxRTHJta6",
				Phone:    "123456",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrorUserNotFound)
				return repo
			},
			ctx: context.Background(),
			user: domain.User{
				Email: "123@qq.com",
				//密码应该是加密的
				Password: "123",
				Phone:    "123456",
				Ctime:    time.Now(),
			},
			wantUser: domain.User{},
			wantErr:  ErrorInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("数据库错误"))
				return repo
			},
			ctx: context.Background(),
			user: domain.User{
				Email: "123@qq.com",
				//密码应该是加密的
				Password: "123",
				Phone:    "123456",
				Ctime:    time.Now(),
			},
			wantUser: domain.User{},
			wantErr:  errors.New("数据库错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						//密码应该是加密的
						Password: "$2a$10$AVk680PliN79.CF3C.KFQeric11w2scHc0slHrK9B9WjPxRTHJta6",
						Phone:    "123456",
						Ctime:    now,
					}, nil)
				return repo
			},
			ctx: context.Background(),
			user: domain.User{
				Email: "123@qq.com",
				//密码应该是加密的
				Password: "1234",
				Phone:    "123456",
				Ctime:    time.Now(),
			},
			wantUser: domain.User{},
			wantErr:  ErrorInvalidUserOrPassword,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			//具体的测试代码
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			user, err := svc.Login(tc.ctx, tc.user)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
