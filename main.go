package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"
	config2 "gohub/config"
	"gohub/pkg/config"

	"github.com/gin-gonic/gin"
)

func init() {
	// 加载 config 目录下的配置信息
	config2.Initialize()
}

func main() {

	// 配置初始化
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()

	config.InitConfig(env)

	// new 一个Gin Engine 实例
	r := gin.New()

	// 初始化DB
	bootstrap.SetupDB()

	// 初始化路由绑定
	bootstrap.SetupRoute(r)
	// 运行服务
	err := r.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err.Error())
	}

}
