package svc

import (
	"log"
	"sea-try-go/service/admin/rpc/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {

	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalln("数据库连接失败")
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
