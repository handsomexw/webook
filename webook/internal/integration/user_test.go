package integration

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()

	testCase := []struct {
		name    string
		reqBody string
		//考虑准备的数据
		before func(t *testing.T)
		//验证数据,处理数据，清除数据，保证每次都是全新的
		after    func(t *testing.T)
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				val, err := rdb.GetDel(ctx, "phone_code:login:123").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},

			reqBody:  `{"phone":"123"}`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 0,
				Msg:  "发送成功",
				Data: nil,
			},
		},
		{
			name: "验证码发送频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, err := rdb.Set(ctx, "phone_code:login:123", "123456",
					0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				val, err := rdb.GetDel(ctx, "phone_code:login:123").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
			},

			reqBody:  `{"phone":"123"}`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "验证码发送错误",
				Data: "验证码发送频繁",
			},
		},
		{
			name: "系统错误，没有设置过期时间",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, err := rdb.Set(ctx, "phone_code:login:123", "123456",
					0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				val, err := rdb.GetDel(ctx, "phone_code:login:123").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
			},

			reqBody:  `{"phone":"123"}`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "验证码发送错误",
				Data: "系统错误",
			},
		},
		{
			name: "手机号输入错误",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},

			reqBody:  `{"phone":""}`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 3,
				Msg:  "手机号输入错误",
				Data: nil,
			},
		},
		{
			name: "数据有误",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},

			reqBody:  `{"phone":}`,
			wantCode: 400,
			wantBody: web.Result{
				Code: 2,
				Msg:  "验证系统错误",
				Data: nil,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/user/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			//resp.Header()
			//var webResp map[string]interface{}
			//err = json.Unmarshal(resp.Body.Bytes(), &webResp)
			var webResp web.Result
			err = json.NewDecoder(resp.Body).Decode(&webResp)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, webResp)
			tc.after(t)
		})
	}
}
