package initializer

import (
	dbModel "whu-campus-auth/model/db"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitAdminUser 初始化默认管理员账户
// 在项目启动时调用，自动创建管理员角色和管理员账户
func InitAdminUser(db *gorm.DB) {
	// 1. 检查是否已有管理员角色
	var adminRole dbModel.Role
	result := db.Where("code = ?", "admin").First(&adminRole)

	if result.Error == gorm.ErrRecordNotFound {
		// 管理员角色不存在，创建之
		adminRole = dbModel.Role{
			Name:   "超级管理员",
			Code:   "admin",
			Desc:   "系统超级管理员，拥有所有权限",
			Status: 1,
		}

		if err := db.Create(&adminRole).Error; err != nil {
			zap.L().Error("创建管理员角色失败", zap.Error(err))
			return
		}

		zap.L().Info("管理员角色创建成功", zap.Uint("id", adminRole.ID))
	} else if result.Error != nil {
		zap.L().Error("查询管理员角色失败", zap.Error(result.Error))
		return
	} else {
		zap.L().Info("管理员角色已存在", zap.Uint("id", adminRole.ID))
	}

	// 2. 检查是否已有管理员账户
	var adminUser dbModel.User
	result = db.Where("username = ?", "admin").First(&adminUser)

	if result.Error == gorm.ErrRecordNotFound {
		// 管理员账户不存在，创建之
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			zap.L().Error("生成密码哈希失败", zap.Error(err))
			return
		}

		adminUser = dbModel.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "系统管理员",
			Email:    "admin@system.local",
			Gender:   0, // 未知
			Status:   1,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			zap.L().Error("创建管理员账户失败", zap.Error(err))
			return
		}

		// 3. 为管理员账户分配管理员角色
		if err := db.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			zap.L().Error("分配管理员角色失败", zap.Error(err))
			return
		}

		zap.L().Info("默认管理员账户创建成功",
			zap.Uint("id", adminUser.ID),
			zap.String("username", adminUser.Username),
			zap.Uint("role_id", adminRole.ID))
	} else if result.Error != nil {
		zap.L().Error("查询管理员账户失败", zap.Error(result.Error))
		return
	} else {
		zap.L().Info("管理员账户已存在", zap.Uint("id", adminUser.ID))

		// 检查是否已有管理员角色
		var userRoles []dbModel.Role
		if err := db.Model(&adminUser).Where("code = ?", "admin").Find(&userRoles).Error; err != nil {
			zap.L().Error("查询管理员角色失败", zap.Error(err))
			return
		}

		if len(userRoles) == 0 {
			// 管理员账户没有管理员角色，分配之
			if err := db.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
				zap.L().Error("补充分配管理员角色失败", zap.Error(err))
				return
			}
			zap.L().Info("为管理员账户补充分配角色成功")
		}
	}

	// 3. 初始化默认字典数据（如果还没有）
	InitDictData(db)
}

// CreateDefaultUser 创建默认测试用户（可选）
// 用于开发环境，创建普通用户账户
func CreateDefaultUser(db *gorm.DB) {
	// 检查是否已有测试用户
	var testUser dbModel.User
	result := db.Where("username = ?", "test").First(&testUser)

	if result.Error == gorm.ErrRecordNotFound {
		// 测试用户不存在，创建之
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			zap.L().Error("生成密码哈希失败", zap.Error(err))
			return
		}

		testUser = dbModel.User{
			Username: "test",
			Password: string(hashedPassword),
			Nickname: "测试用户",
			Email:    "test@example.com",
			Gender:   1, // 男
			Status:   1,
		}

		if err := db.Create(&testUser).Error; err != nil {
			zap.L().Error("创建测试用户失败", zap.Error(err))
			return
		}

		// 获取普通用户角色
		var userRole dbModel.Role
		result = db.Where("code = ?", "user").First(&userRole)

		if result.Error == nil {
			// 分配普通用户角色
			if err := db.Model(&testUser).Association("Roles").Append(&userRole); err != nil {
				zap.L().Error("分配用户角色失败", zap.Error(err))
			}
		}

		zap.L().Info("测试用户创建成功",
			zap.Uint("id", testUser.ID),
			zap.String("username", testUser.Username))
	}
}
