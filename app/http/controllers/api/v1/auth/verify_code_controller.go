package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/pkg/captcha"
	"gohub/pkg/logger"
	"net/http"
)

// VerifyCodeController 用户控制器
type VerifyCodeController struct {
	v1.BaseAPIController
}

// ShowCaptcha 显示图片验证码
func (c VerifyCodeController) ShowCaptcha(ctx *gin.Context) {
	//生成验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	logger.LogIf(err)
	//记录错误日志,因为验证码是用户的入口 出错时应该记 error 等级日志

	// 返回给用户
	ctx.JSON(http.StatusOK, gin.H{
		"captcha_id": id,
		"captcha":    b64s,
	})
}
