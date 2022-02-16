package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gohub/bootstrap"
)

func main() {
	// new 一个Gin Engine 实例air
	r := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(r)

	// 运行服务
	err := r.Run(":8000")
	if err != nil {
		fmt.Println(err.Error())
	}

}
