package dao

import (
	"whu-campus-auth/model/db"

	"gorm.io/gorm"
)

type IRoleDAO interface {
	Create(role *db.Role) error
	GetByID(id uint) (*db.Role, error)
	GetByCode(code string) (*db.Role, error)
	Update(role *db.Role) error
	Delete(id uint) error
	GetList(page, pageSize int, name string, status int) ([]db.Role, int64, error)
	GetAll() ([]db.Role, error)
	PreloadMenus(role *db.Role) error
	AssignMenus(roleID uint, menuIDs []uint) error
	GetDB() *gorm.DB
}

type RoleDAO struct {
	BaseDAO
}

func NewRoleDAO(db *gorm.DB) IRoleDAO {
	return &RoleDAO{
		BaseDAO: NewBaseDAO(db),
	}
}

func (dao *RoleDAO) Create(role *db.Role) error {
	return dao.db.Create(role).Error
}

func (dao *RoleDAO) GetByID(id uint) (*db.Role, error) {
	var role db.Role
	err := dao.db.First(&role, id).Error
	return &role, err
}

func (dao *RoleDAO) GetByCode(code string) (*db.Role, error) {
	var role db.Role
	err := dao.db.Where("code = ?", code).First(&role).Error
	return &role, err
}

func (dao *RoleDAO) Update(role *db.Role) error {
	return dao.db.Save(role).Error
}

func (dao *RoleDAO) Delete(id uint) error {
	return dao.db.Delete(&db.Role{}, id).Error
}

func (dao *RoleDAO) GetList(page, pageSize int, name string, status int) ([]db.Role, int64, error) {
	var roles []db.Role
	var total int64

	query := dao.db.Model(&db.Role{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Preload("Menus").Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error
	return roles, total, err
}

func (dao *RoleDAO) GetAll() ([]db.Role, error) {
	var roles []db.Role
	err := dao.db.Find(&roles).Error
	return roles, err
}

func (dao *RoleDAO) PreloadMenus(role *db.Role) error {
	return dao.db.Preload("Menus").First(role, role.ID).Error
}

func (dao *RoleDAO) AssignMenus(roleID uint, menuIDs []uint) error {
	var role db.Role
	if err := dao.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var menus []db.Menu
	if err := dao.db.Find(&menus, menuIDs).Error; err != nil {
		return err
	}

	return dao.db.Model(&role).Association("Menus").Replace(menus)
}
