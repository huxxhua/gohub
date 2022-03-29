// Package requests 处理请求数据和表单验证
package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/response"
)

// ValidatorFunc 验证函数类型
type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

// Validate 暴露调用
func Validate(ctx *gin.Context, obj interface{}, handler ValidatorFunc) bool {

	//1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := ctx.ShouldBind(obj); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		fmt.Println(err.Error())
		return false
	}

	//2. 表单验证
	errs := handler(obj, ctx)
	if len(errs) > 0 {
		response.ValidationError(ctx, errs)
		return false
	}

	return true
}

func validate(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {

	// 配置初始化
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	// 验证
	return govalidator.New(opts).ValidateStruct()
}

func validateFile(c *gin.Context, data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Request:       c.Request,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}
	// 调用 govalidator 的 Validate 方法来验证文件
	return govalidator.New(opts).Validate()
}
