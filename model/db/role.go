package db

import (
	"time"
	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Desc      string         `gorm:"size:200" json:"desc"`
	Status    int            `gorm:"default:1" json:"status"`
	Menus     []Menu         `gorm:"many2many:role_menus;" json:"menus"`
	Users     []User         `gorm:"many2many:user_roles;" json:"users"`
}

func (Role) TableName() string {
	return "sys_role"
}
