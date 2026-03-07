package initializer

import (
	"whu-campus-auth/model/db"
	"whu-campus-auth/utils"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate(database *gorm.DB) {
	database.AutoMigrate(
		&db.User{},
		&db.Role{},
		&db.Menu{},
		&db.Dict{},
		&db.DictItem{},
	)
	utils.LogInfo("数据库表自动迁移完成")
}
