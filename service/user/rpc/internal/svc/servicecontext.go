package svc

import (
	"fmt"
	"log"
	"sea-try-go/service/user/rpc/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("host=%s user=%s password = %s dbname = %s port = %s sslmode = disable",
		c.Postgres.Host, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("数据库连接失败", err)
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
