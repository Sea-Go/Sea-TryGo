package logic

import (
	"context"

	"sea-try-go/service/article/rpc/internal/model"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesLogic {
	return &ListArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticlesLogic) ListArticles(in *__.ListArticlesRequest) (*__.ListArticlesResponse, error) {
	articles, total, err := l.svcCtx.ArticleRepo.List(l.ctx, model.ListArticlesOption{
		Page:          int(in.Page),
		PageSize:      int(in.PageSize),
		SortBy:        in.SortBy,
		Desc:          in.Desc,
		ManualTypeTag: *in.ManualTypeTag,
		SecondaryTag:  *in.SecondaryTag,
		AuthorId:      *in.AuthorId,
		RelatedGameId: *in.RelatedGameId,
	})

	if err != nil {
		l.Logger.Errorf("ListArticles error: %v", err)
		return nil, err
	}

	var pbArticles []*__.Article
	for _, article := range articles {
		pbArticles = append(pbArticles, &__.Article{
			Id:              article.ID,
			Title:           article.Title,
			Brief:           article.Brief,
			MarkdownContent: article.Content,
			CoverImageUrl:   article.CoverImageURL,
			ManualTypeTag:   article.ManualTypeTag,
			SecondaryTags:   article.SecondaryTags,
			AuthorId:        article.AuthorID,
			CreateTime:      article.CreatedAt.UnixMilli(),
			UpdateTime:      article.UpdatedAt.UnixMilli(),
			Status:          __.ArticleStatus(article.Status),
			ViewCount:       article.ViewCount,
			LikeCount:       article.LikeCount,
			CommentCount:    article.CommentCount,
		})
	}

	return &__.ListArticlesResponse{
		Articles: pbArticles,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
