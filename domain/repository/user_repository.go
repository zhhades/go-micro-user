package repository

import (
	"github.com/zhhades/go-micro-user/domain/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	InitTable() error
	FindUserByName(string) (*model.User, error)
	FindUserById(int64) (*model.User, error)
	CreateUser(*model.User) (int64, error)
	DeleteUserByID(int64) error
	UpdateUser(*model.User) error
	FindAll() ([]model.User, error)
}

type UserRepository struct {
	mysqlDb *gorm.DB
}

func (u *UserRepository) InitTable() error {
	return u.mysqlDb.AutoMigrate(&model.User{})
}

func (u *UserRepository) FindUserByName(s string) (*model.User, error) {
	user := &model.User{}
	return user, u.mysqlDb.Where("user_name = ?", s).Find(user).Error
}

func (u *UserRepository) FindUserById(i int64) (*model.User, error) {
	user := &model.User{}
	return user, u.mysqlDb.First(user, i).Find(user).Error
}

func (u *UserRepository) CreateUser(user *model.User) (int64, error) {
	return user.ID, u.mysqlDb.Create(user).Error
}

func (u *UserRepository) DeleteUserByID(i int64) error {
	return u.mysqlDb.Where("id = ?", i).Delete(&model.User{}).Error
}

func (u *UserRepository) UpdateUser(user *model.User) error {
	return u.mysqlDb.Model(user).Updates(*user).Error
}

func (u *UserRepository) FindAll() ([]model.User, error) {
	userAll := make([]model.User, 0)
	return userAll, u.mysqlDb.Find(&userAll).Error
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		mysqlDb: db,
	}
}
