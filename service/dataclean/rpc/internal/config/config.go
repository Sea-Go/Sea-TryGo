package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type AliGreenConf struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
	Endpoint        string
}

type Config struct {
	zrpc.RpcServerConf
	KqConsumerConf kq.KqConf
	AliGreen       AliGreenConf
}
