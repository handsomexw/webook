package cache

import (
	cmdablemocks "basic-go/webook/internal/repository/cache/redismocks"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCase := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) redis.Cmdable
		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := cmdablemocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "123",
			code:  "123456",

			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := cmdablemocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				res.SetErr(errors.New("mock redis error"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "123",
			code:  "123456",

			wantErr: errors.New("mock redis error"),
		},
		{
			name: "验证码发送频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := cmdablemocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "123",
			code:  "123456",

			wantErr: ErrCodeSentTooMany,
		},
		{
			name: "未找到",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := cmdablemocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-3))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "123",
			code:  "123456",

			wantErr: errors.New("查找失败"),
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := cmdablemocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-10))

				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "123",
			code:  "123456",

			wantErr: errors.New("系统错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cache := NewCodeCache(tc.mock(ctrl))
			err := cache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
