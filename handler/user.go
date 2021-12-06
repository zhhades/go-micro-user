package handler

import (
	"context"
	"github.com/zhhades/go-micro-user/domain/model"

	"github.com/zhhades/go-micro-user/domain/service"
	"github.com/zhhades/go-micro-user/proto/user"
)

type User struct {
	UserDataService service.IUserDataService
}

func (u User) Register(ctx context.Context, request *user.UserRegisterRequest, response *user.UserRegisterResponse) error {
	userRegister := &model.User{
		UserName:     request.UserName,
		FirstName:    request.FirstName,
		HashPassword: request.Pwd,
	}
	_, err := u.UserDataService.AddUser(userRegister)
	if err != nil {
		return err
	}
	response.Msg = "success"
	return nil
}

func (u User) Login(ctx context.Context, request *user.UserLoginRequest, response *user.UserLoginResponse) error {
	isOk, err := u.UserDataService.CheckPwd(request.UserName, request.Pwd)
	if err != nil {
		return err
	}
	response.IsSuccess = isOk
	return nil
}

func (u User) GetUserInfo(ctx context.Context, request *user.UserInfoRequest, response *user.UserInfoResponse) error {
	userInfo, err := u.UserDataService.FindUserByName(request.UserName)
	if err != nil {
		return err
	}
	response = UserForRes(userInfo)
	return nil
}

func UserForRes(userModel *model.User) *user.UserInfoResponse {
	return &user.UserInfoResponse{
		UserName:  userModel.UserName,
		FirstName: userModel.FirstName,
		UserId:    userModel.ID,
	}
}
