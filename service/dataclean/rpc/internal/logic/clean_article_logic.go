package logic

import (
	"context"
	"sea-try-go/service/dataclean/rpc/internal/svc"
	"sea-try-go/service/dataclean/rpc/pb/dataclean"

	"github.com/zeromicro/go-zero/core/logx"
)

type CleanArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanArticleLogic {
	return &CleanArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CleanArticle 提供同步的清洗能力
func (l *CleanArticleLogic) CleanArticle(in *dataclean.CleanArticleRequest) (*dataclean.CleanArticleResponse, error) {

	return nil, nil
}
