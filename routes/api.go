// Package routes 注册路由
package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controllers/api/v1/auth"
	"net/http"
)

func RegisterAPIRoutes(r *gin.Engine) {

	// v1版本路由
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"Hello": "gohub!",
			})
		})

		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
		}
	}
}
