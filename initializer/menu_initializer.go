package initializer

import (
	"whu-campus-auth/dao"
	dbModel "whu-campus-auth/model/db"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitMenus 初始化默认菜单
// 在项目启动时调用，自动创建常用菜单
func InitMenus(db *gorm.DB) {
	menuDAO := dao.NewMenuDAO(db)

	// 检查是否已有菜单
	var count int64
	db.Model(&dbModel.Menu{}).Count(&count)
	zap.L().Info("菜单初始化检查", zap.Int64("当前菜单数量", count))
	if count > 0 {
		zap.L().Info("菜单已存在，跳过初始化")
		return
	}

	// 1. 创建一级菜单：系统管理
	systemMenu := &dbModel.Menu{
		Name:      "系统管理",
		Path:      "/system",
		Component: "layout",
		Icon:      "Setting",
		Sort:      1,
		ParentID:  0,
		Type:      1,
		Status:    1,
	}

	if err := menuDAO.Create(systemMenu); err != nil {
		zap.L().Error("创建系统管理菜单失败", zap.Error(err))
		return
	}

	// 2. 创建二级菜单
	menus := []dbModel.Menu{
		{
			Name:      "用户管理",
			Path:      "/user",
			Component: "user/index",
			Icon:      "User",
			Sort:      1,
			ParentID:  systemMenu.ID,
			Type:      1,
			Status:    1,
		},
		{
			Name:      "角色管理",
			Path:      "/role",
			Component: "role/index",
			Icon:      "UserFilled",
			Sort:      2,
			ParentID:  systemMenu.ID,
			Type:      1,
			Status:    1,
		},
		{
			Name:      "菜单管理",
			Path:      "/menu",
			Component: "menu/index",
			Icon:      "Menu",
			Sort:      3,
			ParentID:  systemMenu.ID,
			Type:      1,
			Status:    1,
		},
		{
			Name:      "字典管理",
			Path:      "/dict",
			Component: "dict/index",
			Icon:      "Collection",
			Sort:      4,
			ParentID:  systemMenu.ID,
			Type:      1,
			Status:    1,
		},
	}

	for i := range menus {
		if err := menuDAO.Create(&menus[i]); err != nil {
			zap.L().Error("创建菜单失败", zap.String("name", menus[i].Name), zap.Error(err))
			continue
		}
	}

	zap.L().Info("菜单初始化成功", zap.Int("count", len(menus)+1))

	// 3. 为管理员角色分配所有菜单
	var adminRole dbModel.Role
	result := db.Where("code = ?", "admin").First(&adminRole)
	if result.Error == nil {
		// 获取所有菜单
		var allMenus []dbModel.Menu
		if err := db.Find(&allMenus).Error; err != nil {
			zap.L().Error("获取所有菜单失败", zap.Error(err))
			return
		}

		// 分配菜单
		if err := db.Model(&adminRole).Association("Menus").Append(allMenus); err != nil {
			zap.L().Error("为管理员角色分配菜单失败", zap.Error(err))
			return
		}

		zap.L().Info("为管理员角色分配菜单成功", zap.Int("count", len(allMenus)))
	}
}
