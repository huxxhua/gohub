// Package user 存放用户 Model 相关逻辑
package user

import (
	"gohub/app/models"
	"gohub/pkg/database"
)

// User 用户模型
type User struct {
	models.BaseModel

	Name string `json:"name,omitempty"`
	// JSON 解析器忽略字段
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimestampsField
}

func (userModel *User) Create() {
	database.DB.Create(&userModel)
}
