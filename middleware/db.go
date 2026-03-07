package middleware

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// DBMiddleware 数据库中间件，将 db 连接放入 context
func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DBContextKey, db)
		c.Next()
	}
}

// GetDB 从 context 中获取数据库连接
func GetDB(c *gin.Context) *gorm.DB {
	db, exists := c.Get(DBContextKey)
	if !exists {
		panic("数据库连接未初始化，请确保使用了 DBMiddleware")
	}
	
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		panic("数据库连接类型错误")
	}
	
	return gormDB
}
