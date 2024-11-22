package interfaces

import (
	"context"

	userApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/user"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
)

type UserInterface struct {
	userApi.UnimplementedUserInterfaceServer
	userUc *biz.UserUseCase
}

func NewUserInterface(userUc *biz.UserUseCase) *UserInterface {
	return &UserInterface{
		userUc: userUc,
	}
}

func (u *UserInterface) GetUserEmail(ctx context.Context, req *userApi.UserEmaileReq) (*userApi.UserEmaileResponse, error) {
	email, err := u.userUc.GetUserEmail(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &userApi.UserEmaileResponse{
		Email: email,
	}, nil
}
