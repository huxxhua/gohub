package v1

import (
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/config"
	"gohub/pkg/file"
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

func (ctrl *UsersController) UpdateEmail(ctx *gin.Context) {

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

func (ctrl *UsersController) UpdatePhone(ctx *gin.Context) {

	request := requests.UserUpdatePhoneRequest{}
	if ok := requests.Validate(ctx, &request, requests.UserUpdatePhone); !ok {
		return
	}

	curUser := auth.CurrentUser(ctx)
	curUser.Phone = request.Phone
	rows := curUser.Save()
	if rows > 0 {
		response.Data(ctx, curUser)
	} else {
		response.Abort500(ctx, "更新失败，请稍后尝试~")
	}
}

func (ctrl *UsersController) UpdatePassword(ctx *gin.Context) {

	request := requests.UserUpdatePasswordRequest{}
	if ok := requests.Validate(ctx, &request, requests.UserUpdatePassword); !ok {
		return
	}

	curUser := auth.CurrentUser(ctx)
	// 验证原始密码是否正确
	_, err := auth.Attempt(curUser.Name, request.Password)
	if err != nil {
		// 失败，显示错误提示
		response.Unauthorized(ctx, "原密码不正确")
	} else {
		// 更新密码为新密码
		curUser.Password = request.NewPassword
		curUser.Save()

		response.Success(ctx)
	}
}

func (ctrl *UsersController) UpdateAvatar(ctx *gin.Context) {

	request := requests.UserUpdateAvatarRequest{}
	if ok := requests.Validate(ctx, &request, requests.UserUpdateAvatar); !ok {
		return
	}

	avatar, err := file.SaveUploadAvatar(ctx, request.Avatar)
	if err != nil {
		response.Abort500(ctx, "上传头像失败，请稍后尝试~\"")
		return
	}

	curUser := auth.CurrentUser(ctx)
	curUser.Avatar = config.GetString("app.url") + avatar
	curUser.Save()

	response.Data(ctx, curUser)
}
