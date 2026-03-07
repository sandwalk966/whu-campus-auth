package dao

import (
	"whu-campus-auth/model/db"

	"gorm.io/gorm"
)

type IMenuDAO interface {
	Create(menu *db.Menu) error
	GetByID(id uint) (*db.Menu, error)
	Update(menu *db.Menu) error
	Delete(id uint) error
	GetList(page, pageSize int, name string, status int) ([]db.Menu, int64, error)
	GetAll() ([]db.Menu, error)
	GetTree() ([]db.Menu, error)
	GetByRoleID(roleID uint) ([]db.Menu, error)
	GetDB() *gorm.DB
}

type MenuDAO struct {
	BaseDAO
}

func NewMenuDAO(db *gorm.DB) IMenuDAO {
	return &MenuDAO{
		BaseDAO: NewBaseDAO(db),
	}
}

func (dao *MenuDAO) Create(menu *db.Menu) error {
	return dao.db.Create(menu).Error
}

func (dao *MenuDAO) GetByID(id uint) (*db.Menu, error) {
	var menu db.Menu
	err := dao.db.First(&menu, id).Error
	return &menu, err
}

func (dao *MenuDAO) Update(menu *db.Menu) error {
	return dao.db.Save(menu).Error
}

func (dao *MenuDAO) Delete(id uint) error {
	return dao.db.Delete(&db.Menu{}, id).Error
}

func (dao *MenuDAO) GetList(page, pageSize int, name string, status int) ([]db.Menu, int64, error) {
	var menus []db.Menu
	var total int64

	query := dao.db.Model(&db.Menu{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&menus).Error
	return menus, total, err
}

func (dao *MenuDAO) GetAll() ([]db.Menu, error) {
	var menus []db.Menu
	err := dao.db.Order("sort ASC").Find(&menus).Error
	return menus, err
}

func (dao *MenuDAO) GetTree() ([]db.Menu, error) {
	var menus []db.Menu
	err := dao.db.Order("parent_id ASC, sort ASC").Find(&menus).Error
	return menus, err
}

func (dao *MenuDAO) GetByRoleID(roleID uint) ([]db.Menu, error) {
	var role db.Role
	if err := dao.db.Preload("Menus").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Menus, nil
}
