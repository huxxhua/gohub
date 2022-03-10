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

// LoginByPassword 多种方法登录，支持手机号、email 和用户名
func (c LoginController) LoginByPassword(ctx *gin.Context) {

	//1. 验证表单
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(ctx, &request, requests.LoginByPassword); !ok {
		return
	}

	// 2. 尝试登录
	userModel, err := auth.Attempt(request.LoginID, request.Password)
	if err != nil {
		response.Error(ctx, err, "登录失败")
	} else {
		toke := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
		response.JSON(ctx, gin.H{
			"token": toke,
		})
	}
}
