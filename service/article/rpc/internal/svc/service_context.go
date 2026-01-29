package svc

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green "github.com/alibabacloud-go/green-20220302/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"sea-try-go/service/article/rpc/internal/config"
	"sea-try-go/service/article/rpc/internal/model"
	"sea-try-go/service/common/snowflake"
)

type ServiceContext struct {
	Config      config.Config
	ArticleRepo *model.ArticleRepo
	KqPusher    *kq.Pusher
	GreenClient *green.Client
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config, articleRepo *model.ArticleRepo) *ServiceContext {

	minioClient, err := minio.New(c.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinIO.AccessKeyID, c.MinIO.SecretAccessKey, ""),
		Secure: c.MinIO.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	err = minioClient.MakeBucket(context.Background(), c.MinIO.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), c.MinIO.BucketName)
		if errBucketExists == nil && exists {
		} else {
			log.Println("Error creating bucket:", err)
		}
	} else {
		policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + c.MinIO.BucketName + `/*"],"Sid": ""}]}`
		err = minioClient.SetBucketPolicy(context.Background(), c.MinIO.BucketName, policy)
		if err != nil {
			log.Println("Error setting bucket policy:", err)
		}
	}

	snowflake.Init()
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
		ArticleRepo: articleRepo,
		KqPusher:    kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		GreenClient: client,
		MinioClient: minioClient,
	}
}
