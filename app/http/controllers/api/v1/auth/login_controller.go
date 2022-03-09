package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

// LoginController 用户控制器
type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 手机登录
func (c *LoginController) LoginByPhone(ctx *gin.Context) {

	// 1. 验证表单
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(ctx, &request, requests.LoginByPhone); !ok {
		return
	}

	// 2. 尝试登录
	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		// 失败，显示错误提示
		response.Error(ctx, err, "账号不存在或密码错误")
	} else {
		// 登录成功
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)

		response.JSON(ctx, gin.H{
			"token": token,
		})
	}

}
