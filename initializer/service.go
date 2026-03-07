package initializer

import "whu-campus-auth/service"

// initService 初始化所有 Service
func (d *Dependencies) initService() {
	d.UserService = service.NewUserService(d.UserDAO)
	d.RoleService = service.NewRoleService(d.RoleDAO)
	d.MenuService = service.NewMenuService(d.MenuDAO)
	d.DictService = service.NewDictService(d.DictDAO)
	d.UploadService = service.NewUploadService()
}
