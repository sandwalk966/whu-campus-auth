package dao

import "gorm.io/gorm"

type BaseDAO struct {
	db *gorm.DB
}

func NewBaseDAO(db *gorm.DB) BaseDAO {
	return BaseDAO{db: db}
}

func (dao *BaseDAO) GetDB() *gorm.DB {
	return dao.db
}
