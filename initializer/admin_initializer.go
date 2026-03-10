package initializer

import (
	"fmt"
	"os"
	dbModel "whu-campus-auth/model/db"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitAdminUser 初始化默认管理员账户
// 在项目启动时调用，自动创建管理员角色和管理员账户
func InitAdminUser(db *gorm.DB) error {
	// 1. 检查是否已有管理员角色
	var adminRole dbModel.Role
	result := db.Where("code = ?", "admin").First(&adminRole)

	if result.Error == gorm.ErrRecordNotFound {
		// 管理员角色不存在，创建之
		adminRole = dbModel.Role{
			Name:   "Super Administrator",
			Code:   "admin",
			Desc:   "System super administrator with all permissions",
			Status: 1,
		}

		if err := db.Create(&adminRole).Error; err != nil {
			zap.L().Error("创建管理员角色失败", zap.Error(err))
			return fmt.Errorf("创建管理员角色失败：%w", err)
		}

		zap.L().Info("管理员角色创建成功", zap.Uint("id", adminRole.ID))
	} else if result.Error != nil {
		zap.L().Error("查询管理员角色失败", zap.Error(result.Error))
		return fmt.Errorf("查询管理员角色失败：%w", result.Error)
	} else {
		zap.L().Info("管理员角色已存在", zap.Uint("id", adminRole.ID))
	}

	// 2. 检查是否已有管理员账户
	var adminUser dbModel.User
	result = db.Where("username = ?", "admin").First(&adminUser)

	if result.Error == gorm.ErrRecordNotFound {
		// 管理员账户不存在，创建之
		password := os.Getenv("ADMIN_PASSWORD")
		if password == "" {
			password = "admin123"
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			zap.L().Error("生成密码哈希失败", zap.Error(err))
			return fmt.Errorf("生成密码哈希失败：%w", err)
		}

		adminUser = dbModel.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "System Administrator",
			Email:    "admin@system.local",
			Gender:   0, // Unknown
			Status:   1,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			zap.L().Error("创建管理员账户失败", zap.Error(err))
			return fmt.Errorf("创建管理员账户失败：%w", err)
		}

		// 3. 为管理员账户分配管理员角色
		if err := db.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			zap.L().Error("分配管理员角色失败", zap.Error(err))
			return fmt.Errorf("分配管理员角色失败：%w", err)
		}

		zap.L().Info("默认管理员账户创建成功",
			zap.Uint("id", adminUser.ID),
			zap.String("username", adminUser.Username),
			zap.Uint("role_id", adminRole.ID))
	} else if result.Error != nil {
		zap.L().Error("查询管理员账户失败", zap.Error(result.Error))
		return fmt.Errorf("查询管理员账户失败：%w", result.Error)
	} else {
		zap.L().Info("管理员账户已存在", zap.Uint("id", adminUser.ID))

		// 检查是否已有管理员角色
		var userRoles []dbModel.Role
		// 重要：需要重新查询用户，因为之前的 adminUser 可能没有完整的 ID 信息
		if err := db.Model(&dbModel.User{}).Where("id = ?", adminUser.ID).Association("Roles").Find(&userRoles); err != nil {
			zap.L().Error("查询用户角色失败", zap.Error(err))
			return fmt.Errorf("查询用户角色失败：%w", err)
		}

		if len(userRoles) == 0 {
			// 管理员账户没有管理员角色，分配之
			if err := db.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
				zap.L().Error("补充分配管理员角色失败", zap.Error(err))
				return fmt.Errorf("补充分配管理员角色失败：%w", err)
			}
			zap.L().Info("为管理员账户补充分配角色成功")
		}
	}

	// 3. 初始化默认字典数据（如果还没有）
	// TODO: 临时注释，测试用
	// InitDictData(db)
	
	return nil
}
