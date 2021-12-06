package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracingV2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/zhhades/go-micro-user/common"
	"github.com/zhhades/go-micro-user/proto/user"
)

func main() {
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"ddns.zhhades.top:8500",
		}
	})

	t, io, err := common.NewTracer("go.micro.service.user.client", "ddns.zhhades.top:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	srv := micro.NewService(
		micro.Name("go.micro.service.user.client"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(consulRegistry),
		micro.WrapHandler(opentracingV2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	userService := user.NewUserService("go.micro.service.user", srv.Client())

	userAdd := &user.UserRegisterRequest{
		UserName:  "zhhades",
		FirstName: "zhhades",
		Pwd:       "zhhades",
	}

	response, err := userService.Register(context.TODO(), userAdd)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)

}
