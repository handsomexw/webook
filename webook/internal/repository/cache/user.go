package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrKeyNotExit = redis.Nil
)

type UserCache struct {
	//面向接口编程
	client     redis.Cmdable
	expiration time.Duration
}

// 依赖注入
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 30,
	}

}

func (u *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := u.key(id)
	result, err := u.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(result, &user)
	return user, err
}

func (u *UserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := u.key(user.Id)
	return u.client.Set(ctx, key, string(val), u.expiration).Err()
}

func (u *UserCache) key(id int64) string {
	return fmt.Sprintf("user:%d", id)
}
