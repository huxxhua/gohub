package v1

import (
	"gohub/app/models/user"
	"gohub/app/requests"
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

	request := requests.PaginationRequest{}
	if ok := requests.Validate(ctx, &request, requests.Pagination); !ok {
		return
	}

	data, pager := user.Paginate(ctx, 10)
	response.JSON(ctx, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (ctrl *UsersController) UpdateProfile(ctx *gin.Context) {

	request := requests.UserUpdateProfileRequest{}
	if ok := requests.Validate(ctx, &request, requests.UserUpdateProfile); !ok {
		return
	}

	curUser := auth.CurrentUser(ctx)
	curUser.Name = request.Name
	curUser.City = request.City
	curUser.Introduction = request.Introduction
	rows := curUser.Save()
	if rows > 0 {
		response.Data(ctx, curUser)
	} else {
		response.Abort500(ctx, "更新失败，请稍后尝试~")
	}
}

func (ctrl UsersController) UpdateEmail(ctx *gin.Context) {

	request := requests.UserUpdateEmailRequest{}
	if ok := requests.Validate(ctx, &request, requests.UserUpdateEmail); !ok {
		return
	}

	curUser := auth.CurrentUser(ctx)
	curUser.Email = request.Email
	rows := curUser.Save()
	if rows > 0 {
		response.Data(ctx, curUser)
	} else {
		response.Abort500(ctx, "更新失败，请稍后尝试~")
	}
}
