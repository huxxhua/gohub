// Package auth 处理用户注册、登录、密码重置
package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/response"
)

// PasswordController 用户控制器
type PasswordController struct {
	v1.BaseAPIController
}

// ResetByPhone 使用手机和验证码重置密码
func (pc *PasswordController) ResetByPhone(ctx *gin.Context) {

	// 1. 验证表单
	request := requests.ResetByPhoneRequest{}

	if ok := requests.Validate(ctx, &request, requests.ResetByPhone); !ok {
		return
	}

	// 2. 更新密码
	userModel := user.GetByPhone(request.Phone)
	if userModel.ID == 0 {
		response.Abort404(ctx)
	} else {
		userModel.Password = request.Password
		userModel.Save()

		response.Success(ctx)
	}
}
