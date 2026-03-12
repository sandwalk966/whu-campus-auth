package initializer

import (
	"fmt"
	"whu-campus-auth/dao"
	dbModel "whu-campus-auth/model/db"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitMenus 初始化默认菜单
// 在项目启动时调用，自动创建常用菜单并分配给管理员角色
func InitMenus(db *gorm.DB) error {
	// zap.L().Info("InitMenus 函数开始执行")
	menuDAO := dao.NewMenuDAO(db)

	// 检查是否已有菜单
	var count int64
	// zap.L().Info("准备检查菜单数量")
	if err := db.Model(&dbModel.Menu{}).Count(&count).Error; err != nil {
		// zap.L().Error("检查菜单数量失败", zap.Error(err))
		return fmt.Errorf("检查菜单数量失败：%w", err)
	}
	// zap.L().Info("菜单初始化检查", zap.Int64("当前菜单数量", count))

	// 如果菜单不存在，创建菜单
	if count == 0 {
		// 1. 创建一级菜单：系统管理
		systemMenu := &dbModel.Menu{
			Name:      "System Management",
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
			return fmt.Errorf("创建系统管理菜单失败：%w", err)
		}

		// 2. 创建二级菜单
		menus := []dbModel.Menu{
			{
				Name:      "User Management",
				Path:      "/user",
				Component: "user/index",
				Icon:      "User",
				Sort:      1,
				ParentID:  systemMenu.ID,
				Type:      1,
				Status:    1,
			},
			{
				Name:      "Role Management",
				Path:      "/role",
				Component: "role/index",
				Icon:      "UserFilled",
				Sort:      2,
				ParentID:  systemMenu.ID,
				Type:      1,
				Status:    1,
			},
			{
				Name:      "Menu Management",
				Path:      "/menu",
				Component: "menu/index",
				Icon:      "Menu",
				Sort:      3,
				ParentID:  systemMenu.ID,
				Type:      1,
				Status:    1,
			},
			{
				Name:      "Dictionary Management",
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
	}

	// 3. 为管理员角色分配所有菜单（无论菜单是否已存在）
	zap.L().Info("开始查询管理员角色")
	var adminRole dbModel.Role
	result := db.Where("code = ?", "admin").First(&adminRole)
	zap.L().Info("管理员角色查询完成", zap.Error(result.Error), zap.Uint("role_id", adminRole.ID))
	if result.Error == nil {
		zap.L().Info("管理员角色已找到", zap.Uint("id", adminRole.ID))

		// 检查管理员角色是否已有菜单权限 - 通过查询关联表
		var menuCount int64
		if err := db.Table("role_menus").Where("role_id = ?", adminRole.ID).Count(&menuCount).Error; err != nil {
			zap.L().Error("检查管理员角色菜单权限失败", zap.Error(err))
			return fmt.Errorf("检查管理员角色菜单权限失败：%w", err)
		}
		zap.L().Info("菜单权限检查结果", zap.Uint("role_id", adminRole.ID), zap.Int64("menu_count", menuCount))

		if menuCount == 0 {
			// 获取所有菜单
			var allMenus []dbModel.Menu
			if err := db.Find(&allMenus).Error; err != nil {
				zap.L().Error("获取所有菜单失败", zap.Error(err))
				return fmt.Errorf("获取所有菜单失败：%w", err)
			}

			// 使用 SQL 直接插入关联关系（避免 GORM 死锁）
			for _, menu := range allMenus {
				if err := db.Exec("INSERT IGNORE INTO role_menus (role_id, menu_id) VALUES (?, ?)", adminRole.ID, menu.ID).Error; err != nil {
					zap.L().Error("插入菜单关联失败", zap.Error(err), zap.Uint("menu_id", menu.ID))
					return fmt.Errorf("插入菜单关联失败：%w", err)
				}
			}

			zap.L().Info("为管理员角色分配菜单成功", zap.Int("count", len(allMenus)))
		} else {
			zap.L().Info("管理员角色已有菜单权限", zap.Int64("count", menuCount))
		}
	} else if result.Error == gorm.ErrRecordNotFound {
		zap.L().Warn("管理员角色不存在，跳过菜单分配")
	} else {
		zap.L().Error("查询管理员角色失败", zap.Error(result.Error))
		return fmt.Errorf("查询管理员角色失败：%w", result.Error)
	}

	return nil
}
