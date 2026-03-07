package initializer

import (
	"gorm.io/gorm"
	"whu-campus-auth/api"
)

// initAPI 初始化所有 API
func (d *Dependencies) initAPI() {
	d.UserAPI = api.NewUserAPI(d.UserService, d.UploadService)
	d.RoleAPI = api.NewRoleAPI(d.RoleService)
	d.MenuAPI = api.NewMenuAPI(d.MenuService)
	d.DictAPI = api.NewDictAPI(d.DictService)
	d.UploadAPI = api.NewUploadAPI(d.UploadService)
}

// InitDependencies 初始化所有依赖
func InitDependencies(db *gorm.DB) *Dependencies {
	deps := &Dependencies{
		DB: db,
	}

	// 按顺序初始化
	deps.initDAO()
	deps.initService()
	deps.initAPI()

	return deps
}
