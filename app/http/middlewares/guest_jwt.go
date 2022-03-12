package middlewares

import (
	"github.com/gin-gonic/gin"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

// GuestJWT 强制使用游客身份访问
func GuestJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.GetHeader("Authorization")) > 0 {

			// 解析 token 成功，说明登录成功了
			_, err := jwt.NewJWT().ParserToken(ctx)
			if err != nil {
				response.Unauthorized(ctx, "请使用游客身份访问")
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
