// Package database 数据库操作
package database

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var SQLDB *sql.DB

// Connect 连接数据库
func Connect(dialector gorm.Dialector, p logger.Interface) {

	// 使用 gorm.Open 连接数据库
	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: p,
	})
	// 处理错误
	if err != nil {
		fmt.Println(err.Error())
	}
	// 获取底层的 sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}
