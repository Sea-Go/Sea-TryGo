package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sea-try-go/service/dataclean/rpc/internal/mqs"
	"sea-try-go/service/dataclean/rpc/pb/dataclean"

	"sea-try-go/service/dataclean/rpc/internal/config"
	"sea-try-go/service/dataclean/rpc/internal/server"
	"sea-try-go/service/dataclean/rpc/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dataclean.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if c.AliGreen.AccessKeyId == "" {
		c.AliGreen.AccessKeyId = os.Getenv("ALI_GREEN_ACCESS_KEY_ID")
	}
	if c.AliGreen.AccessKeySecret == "" {
		c.AliGreen.AccessKeySecret = os.Getenv("ALI_GREEN_ACCESS_KEY_SECRET")
	}

	// 3. 校验密钥是否配置（可选，防止程序启动后报错）
	if c.AliGreen.AccessKeyId == "" || c.AliGreen.AccessKeySecret == "" {
		panic("环境变量 ALI_GREEN_ACCESS_KEY_ID 或 ALI_GREEN_ACCESS_KEY_SECRET 未配置")
	}

	ctx := svc.NewServiceContext(c)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dataclean.RegisterDatacleanServer(grpcServer, server.NewDatacleanServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	serviceGroup.Add(s)

	// Add Kafka consumer
	serviceGroup.Add(kq.MustNewQueue(c.KqConsumerConf, mqs.NewArticleConsumer(context.Background(), ctx)))

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting kafka consumer...\n")
	serviceGroup.Start()
}
