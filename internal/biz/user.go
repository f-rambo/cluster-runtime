package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type UserUseCase struct {
	log *log.Helper
}

func NewUserUseCase(logger log.Logger) *UserUseCase {
	return &UserUseCase{
		log: log.NewHelper(logger),
	}
}

func (u *UserUseCase) GetUserEmail(ctx context.Context, token string) (string, error) {
	githubUser, err := NewClient(token).GetCurrentUser(ctx)
	if err != nil {
		return "", err
	}
	if githubUser == nil {
		return "", errors.New("github user is null")
	}
	return githubUser.GetEmail(), nil
}
