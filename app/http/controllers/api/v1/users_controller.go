package v1

import (
	"gohub/app/models/user"
	"gohub/pkg/auth"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	BaseAPIController
}

func (ctrl *UsersController) CurrentUser(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	response.Data(c, userModel)
}

// Index 所有用户
func (ctrl *UsersController) Index(ctx *gin.Context) {
	data, pager := user.Paginate(ctx, 10)

	response.JSON(ctx, gin.H{
		"data":  data,
		"pager": pager,
	})
}
