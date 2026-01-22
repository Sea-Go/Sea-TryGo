package svc

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green "github.com/alibabacloud-go/green-20220302/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-zero/core/logx"
	"sea-try-go/service/dataclean/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	GreenClient *green.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	config := &openapi.Config{
		AccessKeyId:     &c.AliGreen.AccessKeyId,
		AccessKeySecret: &c.AliGreen.AccessKeySecret,
		Endpoint:        tea.String(c.AliGreen.Endpoint),
		ConnectTimeout:  tea.Int(3000),
		ReadTimeout:     tea.Int(6000),
	}
	client, err := green.NewClient(config)
	if err != nil {
		logx.Errorf("Failed to init AliGreen client: %v", err)
	}

	return &ServiceContext{
		Config:      c,
		GreenClient: client,
	}
}
