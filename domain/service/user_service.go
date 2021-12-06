package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/zhhades/go-micro-user/domain/model"
	"github.com/zhhades/go-micro-user/domain/repository"
)

type IUserDataService interface {
	AddUser(*model.User) (int64, error)
	DeleteUser(int64) error
	UpdateUser(user *model.User, isChangePwd bool) error
	CheckPwd(userName string, pwd string) (isOk bool, err error)
	FindUserByName(string) (*model.User, error)
}

type UserDataService struct {
	UserRepository repository.IUserRepository
}

func NewUserDataService(userRepository repository.IUserRepository) IUserDataService {
	return &UserDataService{
		UserRepository: userRepository,
	}
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(passwd string, hashed string) (isOk bool, err error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(passwd)); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (u *UserDataService) AddUser(user *model.User) (int64, error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return user.ID, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.CreateUser(user)
}

func (u *UserDataService) DeleteUser(i int64) error {
	return u.UserRepository.DeleteUserByID(i)

}

func (u *UserDataService) UpdateUser(user *model.User, isChangePwd bool) error {
	if isChangePwd {
		pwdByte, err := GeneratePassword(user.HashPassword)
		if err != nil {
			return err
		}
		user.HashPassword = string(pwdByte)
	}

	return u.UserRepository.UpdateUser(user)
}

func (u *UserDataService) CheckPwd(userName string, pwd string) (isOk bool, err error) {
	user, err := u.UserRepository.FindUserByName(userName)
	if err != nil {
		return false, err
	}
	return ValidatePassword(pwd, user.HashPassword)
}

func (u *UserDataService) FindUserByName(s string) (*model.User, error) {
	return u.UserRepository.FindUserByName(s)
}
