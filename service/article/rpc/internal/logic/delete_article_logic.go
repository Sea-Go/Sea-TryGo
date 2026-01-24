package logic

import (
	"context"

	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteArticleLogic) DeleteArticle(in *__.DeleteArticleRequest) (*__.DeleteArticleResponse, error) {
	if err := l.svcCtx.ArticleRepo.Delete(l.ctx, in.ArticleId); err != nil {
		l.Logger.Errorf("DeleteArticle db error: %v", err)
		return nil, err
	}

	return &__.DeleteArticleResponse{
		Success: true,
	}, nil
}
