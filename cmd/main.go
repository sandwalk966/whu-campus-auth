package main

import (
	"fmt"
	"whu-campus-auth/config"
	"whu-campus-auth/initializer"
	"whu-campus-auth/middleware"
	"whu-campus-auth/router"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量文件（可选）
	if err := godotenv.Load(".env"); err != nil {
		// .env 文件不存在也没关系，使用默认配置
		fmt.Println("未找到 .env 文件，使用默认配置")
	}

	// 加载配置文件
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("加载配置文件失败：%v", err))
	}

	// 使用环境变量覆盖配置文件
	config.ApplyEnvOverrides(cfg)

	// 初始化日志
	if err := initializer.InitLogger(cfg); err != nil {
		panic(fmt.Sprintf("初始化日志失败：%v", err))
	}
	defer initializer.SyncLogger()

	// 初始化数据库
	db := initializer.InitDatabase(cfg)
	initializer.AutoMigrate(db)

	// 初始化默认管理员账户（先于字典和菜单初始化）
	if err := initializer.InitAdminUser(db); err != nil {
		initializer.LogFatalf("初始化管理员账户失败：%v", err)
	}

	// 初始化默认菜单
	if err := initializer.InitMenus(db); err != nil {
		initializer.LogFatalf("初始化菜单失败：%v", err)
	}

	// 初始化全局数据库连接（供 middleware 使用）
	middleware.InitDB(db)

	// 初始化 Redis
	if err := initializer.InitRedis(&cfg.Redis); err != nil {
		initializer.LogErrorf("Redis 连接失败：%v", err)
	}
	defer initializer.CloseRedis()

	// 初始化所有依赖（DAO → Service → API）
	deps := initializer.InitDependencies(db)

	// 初始化路由并启动服务
	r := initializer.InitRouter(deps)
	router.RegisterStaticRoutes(r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	initializer.LogInfof("服务器启动成功，监听地址：%s", addr)
	if err := r.Run(addr); err != nil {
		initializer.LogFatalf("启动服务器失败：%v", err)
	}
}
