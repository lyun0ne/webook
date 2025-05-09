package service

import (
	"context"
	"errors"

	"github.com/lyun0ne/webook/internal/domain"
	"github.com/lyun0ne/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("邮箱或密码错误")
	ErrUserNotFound          = repository.ErrUserNotFound
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, loginu domain.User) (domain.User, error) {
	findu, err := svc.repo.FindByEmail(ctx, loginu.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(findu.Password), []byte(loginu.Password))
	if err != nil {
		//考虑DEBUG 打印日志
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return findu, nil
}
