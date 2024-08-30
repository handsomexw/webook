package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCase := []struct {
		name string
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		user User

		wantErr error
		wantId  int64
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				res := sqlmock.NewResult(3, 1)
				//gorm 初始化表会加s
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(res)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
				},
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				//res := sqlmock.NewResult(3, 1)
				//gorm 初始化表会加s
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			wantErr: ErrorUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				require.NoError(t, err)
				//res := sqlmock.NewResult(3, 1)
				//gorm 初始化表会加s
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("数据库错误"))
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormmysql.New(gormmysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			d := NewUserDao(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
			//assert.Equal(t, tc.wantId, tc.user.Id)
		})

	}
}
