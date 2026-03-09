package initializer

import "whu-campus-auth/service"

// initService 初始化所有 Service
func (d *Dependencies) initService() {
	// 初始化 Redis 服务
	redisService := service.NewRedisService()
	
	// 初始化所有 Service，注入 Redis 服务
	d.UserService = service.NewUserService(d.UserDAO, redisService)
	d.RoleService = service.NewRoleService(d.RoleDAO)
	d.MenuService = service.NewMenuService(d.MenuDAO)
	d.DictService = service.NewDictService(d.DictDAO)
	d.UploadService = service.NewUploadService()
}
