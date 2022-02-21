// Package auth 处理用户身份认证相关逻辑
package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"net/http"
)

// SignupController 注册控制器
type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist 检测手机号是否被注册
func (c *SignupController) IsPhoneExist(ctx *gin.Context) {

	// 初始化请求对象
	request := requests.SignupPhoneExistRequest{}

	// 解释JSON请求
	if err := ctx.ShouldBindJSON(&request); err != nil {
		// 解释失败 返回422 状态码和错误信息
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		// 打印错误信息
		fmt.Println(err.Error())
		return
	}
	// 表单验证
	errs := requests.ValidateSignupPhoneExist(&request, ctx)
	// errs 返回长度等于零即通过，大于 0 即有错误发生
	if len(errs) > 0 {
		// 验证失败，返回 422 状态码和错误信息
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errs,
		})
		return
	}

	//  检查数据库并返回响应
	ctx.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

// IsEmailExist 检测邮箱是否已注册
func (c *SignupController) IsEmailExist(ctx *gin.Context) {

	// 初始化请求对象
	request := requests.SignupEmailExistRequest{}

	// 解析 JSON 请求
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		// 打印错误信息
		fmt.Println(err.Error())
		return
	}
	// 表单验证
	errs := requests.ValidateSignupEmailExist(&request, ctx)
	if len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errs,
		})
		return
	}

	//  检查数据库并返回响应
	ctx.JSON(http.StatusOK, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}
