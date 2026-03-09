package dao

import (
	"whu-campus-auth/model/db"

	"gorm.io/gorm"
)

type IUserDAO interface {
	Create(user *db.User) error
	GetByID(id uint) (*db.User, error)
	GetByUsername(username string) (*db.User, error)
	Update(user *db.User) error
	Delete(id uint) error
	GetList(page, pageSize int, username string, status int) ([]db.User, int64, error)
	PreloadRoles(user *db.User) error
	AssignRoles(userID uint, roleIDs []uint) error
	GetDB() *gorm.DB
}

type UserDAO struct {
	BaseDAO
}

func NewUserDAO(db *gorm.DB) IUserDAO {
	return &UserDAO{
		BaseDAO: NewBaseDAO(db),
	}
}

func (dao *UserDAO) Create(user *db.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDAO) GetByID(id uint) (*db.User, error) {
	var user db.User
	err := dao.db.First(&user, id).Error
	return &user, err
}

func (dao *UserDAO) GetByUsername(username string) (*db.User, error) {
	var user db.User
	err := dao.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (dao *UserDAO) Update(user *db.User) error {
	return dao.db.Save(user).Error
}

func (dao *UserDAO) Delete(id uint) error {
	return dao.db.Delete(&db.User{}, id).Error
}

func (dao *UserDAO) GetList(page, pageSize int, username string, status int) ([]db.User, int64, error) {
	var users []db.User
	var total int64

	query := dao.db.Model(&db.User{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Preload("Roles").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (dao *UserDAO) PreloadRoles(user *db.User) error {
	return dao.db.Preload("Roles.Menus").First(user, user.ID).Error
}

func (dao *UserDAO) AssignRoles(userID uint, roleIDs []uint) error {
	var user db.User
	if err := dao.db.First(&user, userID).Error; err != nil {
		return err
	}

	var roles []db.Role
	if err := dao.db.Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	return dao.db.Model(&user).Association("Roles").Replace(roles)
}
