package mysql

import (
	"zero-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// NewDB 获取数据库连接
func NewDB(cfg *config.Config, logger logger.Interface) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.MySQL.Dsn), &gorm.Config{
		Logger: logger,

		// 在自动迁移时，忽略外键约束
		DisableForeignKeyConstraintWhenMigrating: true,

		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.MySQL.Prefix,
			SingularTable: true,
		},
	})

	if err != nil {
		panic(err)
	}

	return db
}
