package initializer

import (
	"gorm.io/gorm"
	"whu-campus-auth/api"
	"whu-campus-auth/dao"
	"whu-campus-auth/service"
)

// Dependencies 依赖容器
type Dependencies struct {
	DB *gorm.DB

	// DAO 层
	UserDAO  dao.IUserDAO
	RoleDAO  dao.IRoleDAO
	MenuDAO  dao.IMenuDAO
	DictDAO  dao.IDictDAO

	// Service 层
	UserService    *service.UserService
	RoleService    *service.RoleService
	MenuService    *service.MenuService
	DictService    *service.DictService
	UploadService  *service.UploadService

	// API 层
	UserAPI    *api.UserAPI
	RoleAPI    *api.RoleAPI
	MenuAPI    *api.MenuAPI
	DictAPI    *api.DictAPI
	UploadAPI  *api.UploadAPI
}
