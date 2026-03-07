package db

import (
	"time"
	"gorm.io/gorm"
)

type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Path      string         `gorm:"size:200" json:"path"`
	Component string         `gorm:"size:200" json:"component"`
	Icon      string         `gorm:"size:50" json:"icon"`
	Sort      int            `gorm:"default:0" json:"sort"`
	ParentID  uint           `gorm:"default:0" json:"parent_id"`
	Type      int            `gorm:"default:1" json:"type"`
	Status    int            `gorm:"default:1" json:"status"`
	Children  []Menu         `gorm:"foreignKey:ParentID" json:"children"`
	Roles     []Role         `gorm:"many2many:role_menus;" json:"roles"`
}

func (Menu) TableName() string {
	return "sys_menu"
}
