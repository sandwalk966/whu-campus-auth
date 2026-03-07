package db

import (
	"time"
	"gorm.io/gorm"
)

type Dict struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Desc      string         `gorm:"size:200" json:"desc"`
	Status    int            `gorm:"default:1" json:"status"`
	Items     []DictItem     `json:"items"`
}

func (Dict) TableName() string {
	return "sys_dict"
}

type DictItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DictID    uint           `gorm:"not null;index" json:"dict_id"`
	Label     string         `gorm:"size:50;not null" json:"label"`
	Value     string         `gorm:"size:100;not null" json:"value"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int            `gorm:"default:1" json:"status"`
	Dict      Dict           `json:"-"`
}

func (DictItem) TableName() string {
	return "sys_dict_item"
}
