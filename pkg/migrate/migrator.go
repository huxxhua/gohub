// Package migrate 处理数据库迁移
package migrate

import (
	"gohub/pkg/console"
	"gohub/pkg/database"
	"gohub/pkg/file"
	"gohub/pkg/logger"
	"gorm.io/gorm"
	"io/ioutil"
)

// Migrator 数据迁移操作类
type Migrator struct {
	Folder   string
	DB       *gorm.DB
	Migrator gorm.Migrator
}

// Migration 对应数据的 migrations 表里的一条数据
type Migration struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Migration string `gorm:"type:varchar(255);not null;unique;"`
	Batch     int
}

// NewMigrator 创建 Migrator 实例，用以执行迁移操作
func NewMigrator() *Migrator {

	// 初始化必要属性
	migrator := &Migrator{
		Folder:   "database/migrations/",
		DB:       database.DB,
		Migrator: database.DB.Migrator(),
	}
	// migrations 不存在的话就创建它
	migrator.createMigrationsTable()

	return migrator
}

// 创建 migrations 表
func (m *Migrator) createMigrationsTable() {

	migration := Migration{}

	// 不存在才创建
	if !m.Migrator.HasTable(&migration) {
		err := m.Migrator.CreateTable(&migration)
		if err != nil {
			logger.LogIf(err)
		}
	}
}

// Up 执行所有未迁移过的文件
func (m *Migrator) Up() {

	// 读取所有迁移文件，确保按照时间排序
	migrateFile := m.readAllMigrationFiles()

	// 获取当前批次的值
	batch := m.getBatch()

	// 获取所有迁移数据
	var migrations []Migration
	m.DB.Find(&migrations)

	//m.DB.Where("id = ?",10).Find(&migrations)
	// 可以通过此值来判断数据库是否已是最新
	runed := false

	// 对迁移文件进行遍历，如果没有执行过，就执行 up 回调
	for _, mfile := range migrateFile {

		// 对比文件名称，看是否已经运行过
		if mfile.isNotMigrated(migrations) {
			m.runUpMigration(mfile, batch)
			runed = true
		}
	}

	if !runed {
		console.Success("database is up to date.")
	}
}

// Rollback 回滚上一个操作
func (m *Migrator) Rollback() {

	// 获取最后一批次的迁移数据
	lastMigration := Migration{}
	m.DB.Order("id DESC").First(&lastMigration)

	var migrations []Migration
	m.DB.Where("batch = ?", lastMigration.Batch).Order("id DESC").Find(&migrations)

	// 回滚最后一批次的迁移
	if !m.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to rollback.")
	}
}

// 回退迁移，按照倒序执行迁移的 down 方法
func (m *Migrator) rollbackMigrations(migration []Migration) bool {

	// 标记是否真的有执行了迁移回退的操作
	runed := false

	for _, migrate := range migration {

		// 友好提示
		console.Warning("rollback " + migrate.Migration)

		// 执行迁移文件的 down 方法
		mfile := getMigrationFile(migrate.Migration)
		if mfile.Down != nil {
			mfile.Down(database.DB.Migrator(), database.SQLDB)
		}

		runed = true

		// 回退成功了就删除掉这条记录
		m.DB.Delete(&migrate)

		// 打印运行状态
		console.Success("finish " + mfile.FileName)
	}

	return runed
}

// 获取当前这个批次的值
func (m *Migrator) getBatch() int {

	// 默认为 1
	batch := 1

	// 取最后执行的一条迁移数据
	lastMigration := Migration{}
	m.DB.Order("id DESC").First(&lastMigration)

	// 如果有值的话，加一

	if lastMigration.ID > 0 {
		batch = lastMigration.Batch + 1
	}

	return batch
}

// 从文件目录读取文件，保证正确的时间排序
func (m *Migrator) readAllMigrationFiles() []MigrationFile {

	// 读取 database/migrations/ 目录下的所有文件
	// 默认是会按照文件名称进行排序
	files, err := ioutil.ReadDir(m.Folder)
	console.ExitIf(err)

	var migrateFiles []MigrationFile
	for _, f := range files {

		// 去除文件后缀 .go
		fileName := file.FileNameWithoutExtension(f.Name())

		// 通过迁移文件的名称获取『MigrationFile』对象
		mfile := getMigrationFile(fileName)

		// 加个判断，确保迁移文件可用，再放进 migrateFiles 数组中
		if len(mfile.FileName) > 0 {
			migrateFiles = append(migrateFiles, mfile)
		}
	}

	// 返回排序好的『MigrationFile』数组
	return migrateFiles
}

// 执行迁移，执行迁移的 up 方法
func (m *Migrator) runUpMigration(mfile MigrationFile, batch int) {

	// 执行 up 区块的 SQL
	if mfile.Up != nil {
		// 友好提示
		console.Warning("migrating " + mfile.FileName)
		// 执行 up 方法
		mfile.Up(database.DB.Migrator(), database.SQLDB)
		// 提示已迁移了哪个文件
		console.Success("migrated " + mfile.FileName)
	}

	// 入库
	err := m.DB.Create(&Migration{
		Migration: mfile.FileName,
		Batch:     batch,
	}).Error

	console.ExitIf(err)
}

// Reset 回滚所有迁移
func (m *Migrator) Reset() {

	var migrations []Migration

	// 按照倒序读取所有迁移文件
	m.DB.Order("id DESC").First(&migrations)

	// 回滚所有迁移
	if !m.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to reset.")
	}
}

// Refresh 回滚所有迁移，并运行所有迁移
func (m *Migrator) Refresh() {

	// 回滚所有迁移
	m.Reset()

	// 再次执行所有迁移
	m.Up()
}

// Fresh Drop 所有的表并重新运行所有迁移
func (m *Migrator) Fresh() {

	// 获取数据库名称，用以提示
	dbname := database.CurrentDatabase()

	// 删除所有表
	err := database.DeleteAllTables()
	console.ExitIf(err)
	console.Success("clear up database" + dbname)

	// 重新创建 migrates 表
	m.createMigrationsTable()
	console.Success("[migration] table created")

	// 重新调用 up 命令
	m.Up()
}
