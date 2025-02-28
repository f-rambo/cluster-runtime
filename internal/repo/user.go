package repo

import (
	"context"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type UserRepo struct {
	log *log.Helper
}

func NewUserRepo(logger log.Logger) biz.UserRepoInterface {
	return &UserRepo{log: log.NewHelper(logger)}
}

func (u *UserRepo) GetGithubUserEmail(ctx context.Context, token string) (string, error) {
	githubUser, err := NewClient(token).GetCurrentUser(ctx)
	if err != nil {
		return "", err
	}
	if githubUser == nil {
		return "", errors.New("github user is null")
	}
	return githubUser.GetEmail(), nil
}
