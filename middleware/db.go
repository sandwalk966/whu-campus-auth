package middleware

import (
	"context"
	"whu-campus-auth/config"

	"gorm.io/gorm"
)

var globalDB *gorm.DB

// InitDB 初始化全局数据库连接（在 main.go 中调用）
func InitDB(db *gorm.DB) {
	globalDB = db
}

// GetDB 获取全局数据库连接
func GetDB() *gorm.DB {
	if globalDB == nil {
		cfg := config.GlobalConfig
		if cfg == nil {
			return nil
		}
		// 懒加载初始化
		globalDB = initLazyDB(cfg)
	}
	return globalDB
}

func initLazyDB(cfg *config.Config) *gorm.DB {
	// 这里不应该再次初始化，因为 main.go 已经初始化了
	// 如果到这里，说明程序启动顺序有问题
	return nil
}

// GetDBFromContext 从 context 获取数据库连接（备用方案）
func GetDBFromContext(ctx context.Context) *gorm.DB {
	db, _ := ctx.Value("db").(*gorm.DB)
	return db
}
