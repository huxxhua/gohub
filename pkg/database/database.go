// Package database 数据库操作
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"gohub/pkg/config"
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

func CurrentDatabase() (dbname string) {
	dbname = DB.Migrator().CurrentDatabase()
	return
}

func DeleteAllTables() error {
	var err error
	switch config.Get("database.connection") {
	case "mysql":
		err = deleteAllMySQLTables()
	case "sqlite":
		err = deleteAllSqliteTables()
	default:
		panic(errors.New("database connection not supported"))
	}

	return err
}

func deleteAllSqliteTables() error {

	tables := []string{}

	// 读取所有数据表
	err := DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteAllMySQLTables() error {

	dbname := CurrentDatabase()
	var tables []string

	// 读取所有数据表
	err := DB.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// 暂时关闭外键检测
	DB.Exec("SET foreign_key_checks = 0;")

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// 开启 MySQL 外键检测
	DB.Exec("SET foreign_key_checks = 1;")

	return nil
}

func TableName(obj interface{}) string {
	stmt := &gorm.Statement{DB: DB}
	err := stmt.Parse(obj)
	if err != nil {
		return ""
	}
	return stmt.Schema.Table
}
