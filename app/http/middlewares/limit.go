package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gohub/pkg/app"
	"gohub/pkg/limiter"
	"gohub/pkg/logger"
	"gohub/pkg/response"
	"net/http"
)

// LimitIP 全局限流中间件，针对 IP 进行限流
// limit 为格式化字符串，如 "5-S" ，示例:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"

func LimitIP(limit string) gin.HandlerFunc {

	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(ctx *gin.Context) {
		// 针对 IP 限流
		key := limiter.GetKeyIP(ctx)
		if ok := limitHandler(ctx, key, limit); !ok {
			return
		}
		ctx.Next()
	}
}

// LimitPerRoute 限流中间件，用在单独的路由中
func LimitPerRoute(limit string) gin.HandlerFunc {

	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(ctx *gin.Context) {

		// 针对单个路由，增加访问次数
		ctx.Set("limiter-once", false)

		// 针对 IP + 路由进行限流
		key := limiter.GetKeyRouteWithIP(ctx)
		if ok := limitHandler(ctx, key, limit); !ok {
			return
		}

		ctx.Next()
	}
}

func limitHandler(ctx *gin.Context, key string, limit string) bool {

	// 获取超额的情况
	rate, err := limiter.CheckRate(ctx, key, limit)
	if err != nil {
		logger.LogIf(err)
		response.Abort500(ctx)
		return false
	}

	// ---- 设置标头信息-----
	// X-RateLimit-Limit :10000 最大访问次数
	// X-RateLimit-Remaining :9993 剩余的访问次数
	// X-RateLimit-Reset :1513784506 到该时间点，访问次数会重置为 X-RateLimit-Limit
	ctx.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))
	ctx.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining))
	ctx.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))

	// 超额
	if rate.Reached {
		// 提示用户超额了
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "接口请求太频繁",
		})
		return false
	}

	return true
}
