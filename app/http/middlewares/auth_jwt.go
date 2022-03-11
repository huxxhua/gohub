// Package middlewares Gin 中间件
package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gohub/app/models/user"
	"gohub/pkg/config"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

func AuthJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 从标头 Authorization:Bearer xxxxx 中获取信息，并验证 JWT 的准确性
		claims, err := jwt.NewJWT().ParserToken(ctx)

		// JWT 解析失败，有错误发生
		if err != nil {
			response.Unauthorized(ctx, fmt.Sprintf("请查看 %v 相关接口认证文档", config.GetString("app.name")))
			return
		}

		// JWT 解析成功，设置用户信息
		userModel := user.Get(claims.UserID)
		if userModel.ID == 0 {
			response.Unauthorized(ctx, "找不到对应用户，用户可能已删除")
			return
		}

		// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
		ctx.Set("current_user_id", userModel.GetStringID())
		ctx.Set("current_user_name", userModel.Name)
		ctx.Set("current_user", userModel)

		// 当在 c.Next() 之前 return 掉，就会中断所有的后续请求
		ctx.Next()
	}
}
