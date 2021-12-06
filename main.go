package main

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/zhhades/go-micro-user/domain/repository"
	"github.com/zhhades/go-micro-user/domain/service"
	"github.com/zhhades/go-micro-user/handler"
	"github.com/zhhades/go-micro-user/proto/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
	)

	srv.Init()

	dsn := "root:12345678@tcp(192.168.0.188:3306)/zhhades?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	rp := repository.NewUserRepository(db)
	err = rp.InitTable()
	if err != nil {
		fmt.Println(err)
	}

	userDataService := service.NewUserDataService(rp)

	err = user.RegisterUserHandler(srv.Server(), &handler.User{
		UserDataService: userDataService,
	})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}
