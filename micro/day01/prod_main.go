package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"www.gomicro.com/prodservice"

)

func main() {

	consulReg := consul.NewRegistry(
		registry.Addrs("192.168.86.113:8500"),
		)

	ginRouter := gin.Default()
	v1Group := ginRouter.Group("/v1")

	{
		v1Group.Handle("POST","/prods",func(ctx *gin.Context){
			ctx.JSON(200,gin.H{
				"data": prodservice.NewProdList(5),
			})
		})
	}


	server := web.NewService(
		web.Name("prodservice"),
		web.Address(":8001"),
		web.Handler(ginRouter),
		web.Registry(consulReg),
		)

	server.Run()
}
