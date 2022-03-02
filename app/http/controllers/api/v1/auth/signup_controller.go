// Package auth 处理用户身份认证相关逻辑
package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/response"
)

// SignupController 注册控制器
type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist 检测手机号是否被注册
func (c *SignupController) IsPhoneExist(ctx *gin.Context) {

	// 初始化请求对象
	request := requests.SignupPhoneExistRequest{}

	if err := requests.Validate(ctx, &request, requests.ValidateSignupPhoneExist); err != true {
		return
	}
	//  检查数据库并返回响应
	response.JSON(ctx, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

// IsEmailExist 检测邮箱是否已注册
func (c *SignupController) IsEmailExist(ctx *gin.Context) {

	// 初始化请求对象
	request := requests.SignupEmailExistRequest{}

	if err := requests.Validate(ctx, &request, requests.ValidateSignupEmailExist); err != true {
		return
	}

	//  检查数据库并返回响应
	response.JSON(ctx, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}
