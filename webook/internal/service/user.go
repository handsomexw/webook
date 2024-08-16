package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserDuplicate         = repository.ErrorUserDuplicate
	ErrorInvalidUserOrPassword = errors.New("邮箱/密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)

	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, user.Email)
	if errors.Is(err, repository.ErrorUserNotFound) {
		return domain.User{}, ErrorInvalidUserOrPassword
	}
	//系统错误
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		//打印日志
		return domain.User{}, ErrorInvalidUserOrPassword
	}
	return u, err
	//两个错误：系统错误和密码错误
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, id)

	return user, err
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//这是快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrorUserNotFound {
		return u, err
	}
	//这是慢路径
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != ErrorUserDuplicate {
		return u, err
	}
	//未找到就创建一个后再找一遍就找到了
	//存在主从延迟
	return svc.repo.FindByPhone(ctx, phone)
}
