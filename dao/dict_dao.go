package dao

import (
	"whu-campus-auth/model/db"

	"gorm.io/gorm"
)

type IDictDAO interface {
	Create(dict *db.Dict) error
	GetByID(id uint) (*db.Dict, error)
	GetByCode(code string) (*db.Dict, error)
	Update(dict *db.Dict) error
	Delete(id uint) error
	GetList(page, pageSize int, name string, status int) ([]db.Dict, int64, error)
	PreloadItems(dict *db.Dict) error
	GetDB() *gorm.DB
}

type DictDAO struct {
	BaseDAO
}

func NewDictDAO(db *gorm.DB) IDictDAO {
	return &DictDAO{
		BaseDAO: NewBaseDAO(db),
	}
}

func (dao *DictDAO) Create(dict *db.Dict) error {
	return dao.db.Create(dict).Error
}

func (dao *DictDAO) GetByID(id uint) (*db.Dict, error) {
	var dict db.Dict
	err := dao.db.First(&dict, id).Error
	return &dict, err
}

func (dao *DictDAO) GetByCode(code string) (*db.Dict, error) {
	var dict db.Dict
	err := dao.db.Where("code = ?", code).First(&dict).Error
	return &dict, err
}

func (dao *DictDAO) Update(dict *db.Dict) error {
	return dao.db.Save(dict).Error
}

func (dao *DictDAO) Delete(id uint) error {
	return dao.db.Delete(&db.Dict{}, id).Error
}

func (dao *DictDAO) GetList(page, pageSize int, name string, status int) ([]db.Dict, int64, error) {
	var dicts []db.Dict
	var total int64

	query := dao.db.Model(&db.Dict{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Preload("Items").Offset((page - 1) * pageSize).Limit(pageSize).Find(&dicts).Error
	return dicts, total, err
}

func (dao *DictDAO) PreloadItems(dict *db.Dict) error {
	return dao.db.Preload("Items").First(dict, dict.ID).Error
}
