// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"

	"sea-try-go/service/article/api/internal/svc"
	"sea-try-go/service/article/api/internal/types"
	"sea-try-go/service/article/rpc/articleservice"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleLogic) UpdateArticle(req *types.UpdateArticleReq) (resp *types.UpdateArticleResp, err error) {
	// Construct RPC request, handling optional fields
	rpcReq := &articleservice.UpdateArticleRequest{
		ArticleId:     req.ArticleId,
		SecondaryTags: req.SecondaryTags,
	}

	if req.Title != "" {
		rpcReq.Title = &req.Title
	}
	if req.Brief != "" {
		rpcReq.Brief = &req.Brief
	}
	if req.Content != "" {
		rpcReq.MarkdownContent = &req.Content
	}
	if req.CoverImageUrl != "" {
		rpcReq.CoverImageUrl = &req.CoverImageUrl
	}
	if req.ManualTypeTag != "" {
		rpcReq.ManualTypeTag = &req.ManualTypeTag
	}
	if req.Status != 0 {
		status := __.ArticleStatus(req.Status)
		rpcReq.Status = status.Enum()
	}

	res, err := l.svcCtx.ArticleRpc.UpdateArticle(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.UpdateArticleResp{
		Success: res.Success,
	}, nil
}
