package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func main() {
	// new 一个Gin Engine 实例air
	r := gin.New()

	// 注册中间件
	r.Use(gin.Logger(), gin.Recovery())

	// 注册一个路由
	r.GET("/", func(c *gin.Context) {

		// 以 JSON 格式响应
		c.JSON(http.StatusOK, gin.H{
			"Hello": "gohub!",
		})
	})

	// 处理 404 请求
	r.NoRoute(func(c *gin.Context) {
		// 获取标头信息 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是HTML的话
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// 默认返回JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})

	// 运行服务
	err := r.Run(":8000")
	if err != nil {
		return
	}

}
