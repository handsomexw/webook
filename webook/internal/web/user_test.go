package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	svcmocks "basic-go/webook/internal/service/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {

	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
					Return(nil)
				return usersvc
			},
			reqBody: `{
				"email":"123595@qq.com",
				"confirmPassword":"123456",
				"password":"123456"
			}`,
			wantCode: http.StatusOK,
			//wantBody: `{"message":"注册成功"}`,
			//两种比较json
			wantBody: `{"confirm":"123456","message":"注册成功","password":"123456","user":"123595@qq.com"}`,
			//wantBody: "注册成功",
			//wantBody: `{"code":0,"msg":"注册成功","data":null}`,
		},
		{name: "参数不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				"email":"123595@qq.com",
				"confirmPassword":"123456",
				"password":"123456,
			}`,
			wantCode: http.StatusBadRequest,
		},
		{name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrorUserDuplicate)
				return usersvc
			},
			reqBody: `{
				"email":"123@qq.com",
				"confirmPassword":"123456",
				"password":"123456"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{"message":"邮箱冲突"}`,
		},
		{name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `{
				"email":"123qq.com",
				"confirmPassword":"123456",
				"password":"123456"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{"message":"邮箱格式错误"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			RegisterRoutes(server, h)

			req, err := http.NewRequest(http.MethodPost,
				"/user/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			//resp.Header()
			jsonBody := toString(resp.Body.Bytes())
			tc.wantBody = utoString(tc.wantBody)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, jsonBody)
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)
	//预期行为与预期返回值
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log("err:", err)
}

func toString(body []byte) string {
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return string(body)
	}
	jsonBody, _ := json.Marshal(data)
	return string(jsonBody)
}

func utoString(body string) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return body
	}
	jsonBody, _ := json.Marshal(data)
	return string(jsonBody)
}
