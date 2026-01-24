// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"

	"fmt"

	"sea-try-go/service/article/api/internal/svc"
	"sea-try-go/service/article/api/internal/types"
	"sea-try-go/service/article/rpc/articleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleLogic) GetArticle(req *types.GetArticleReq) (resp *types.GetArticleResp, err error) {
	res, err := l.svcCtx.ArticleRpc.GetArticle(l.ctx, &articleservice.GetArticleRequest{
		ArticleId: req.ArticleId,
		IncrView:  req.IncrView,
	})
	if err != nil {
		return nil, err
	}

	if res.Article == nil {
		return nil, fmt.Errorf("article not found")
	}

	return &types.GetArticleResp{
		Article: types.Article{
			Id:            res.Article.Id,
			Title:         res.Article.Title,
			Brief:         res.Article.Brief,
			Content:       res.Article.MarkdownContent,
			CoverImageUrl: res.Article.CoverImageUrl,
			ManualTypeTag: res.Article.ManualTypeTag,
			SecondaryTags: res.Article.SecondaryTags,
			AuthorId:      res.Article.AuthorId,
			CreateTime:    res.Article.CreateTime,
			UpdateTime:    res.Article.UpdateTime,
			Status:        int32(res.Article.Status),
			ViewCount:     res.Article.ViewCount,
			LikeCount:     res.Article.LikeCount,
			CommentCount:  res.Article.CommentCount,
		},
	}, nil
}
