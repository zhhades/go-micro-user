package main

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracingV2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/zhhades/go-micro-user/common"
	"github.com/zhhades/go-micro-user/domain/repository"
	"github.com/zhhades/go-micro-user/domain/service"
	"github.com/zhhades/go-micro-user/handler"
	"github.com/zhhades/go-micro-user/proto/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	consulConfig, err := common.GetConsulConfig("ddns.zhhades.top", 8500, "/micro/config")

	if err != nil {
		log.Error(err)
	}
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"ddns.zhhades.top:8500",
		}
	})
	t, io, err := common.NewTracer("go.micro.service.user", "ddns.zhhades.top:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(consulRegistry),
		micro.WrapHandler(opentracingV2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	srv.Init()

	mysqlConfig := common.GetMysqlConfigFromConsul(consulConfig, "mysql")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Passwd, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DataBase)

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
