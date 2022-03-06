// Package auth 处理用户身份认证相关逻辑
package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/jwt"
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

// SignupUsingPhone 使用手机和验证码进行注册
func (c *SignupController) SignupUsingPhone(ctx *gin.Context) {

	// 1. 验证表单
	request := requests.SignupUsingPhoneRequest{}
	if ok := requests.Validate(ctx, &request, requests.SignupUsingPhone); !ok {
		return
	}

	// 2. 验证成功，创建数据
	userModel := user.User{
		Name:     request.Name,
		Phone:    request.Phone,
		Password: request.Password,
	}
	userModel.Create()

	if userModel.ID > 0 {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
		response.CreatedJSON(ctx, gin.H{
			"data":  userModel,
			"token": token,
		})
	} else {
		response.Abort500(ctx, "创建用户失败，请稍后尝试~")
	}

}

// SignupUsingEmail 使用 Email + 验证码进行注册
func (c *SignupController) SignupUsingEmail(ctx *gin.Context) {

	// 1. 验证表单
	request := requests.SignupUsingEmailRequest{}
	if ok := requests.Validate(ctx, &request, requests.SignupUsingEmail); !ok {
		return
	}

	// 2. 验证成功，创建数据
	userModel := user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	userModel.Create()

	if userModel.ID > 0 {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
		response.CreatedJSON(ctx, gin.H{
			"data":  userModel,
			"token": token,
		})
	} else {
		response.Abort500(ctx, "创建用户失败，请稍后尝试~")
	}
}
