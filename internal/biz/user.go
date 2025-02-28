package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepoInterface interface {
	GetGithubUserEmail(ctx context.Context, token string) (string, error)
}

type UserUseCase struct {
	repo UserRepoInterface
	log  *log.Helper
}

func NewUserUseCase(repo UserRepoInterface, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (u *UserUseCase) GetUserEmail(ctx context.Context, token string) (string, error) {
	return u.repo.GetGithubUserEmail(ctx, token)
}
