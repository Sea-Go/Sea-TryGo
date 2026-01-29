package logic

import (
	"context"

	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleLogic) GetArticle(in *__.GetArticleRequest) (*__.GetArticleResponse, error) {
	article, err := l.svcCtx.ArticleRepo.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		l.Logger.Errorf("GetArticle db error: %v", err)
		return nil, err
	}

	if in.IncrView {
		if err := l.svcCtx.ArticleRepo.IncrViewCount(l.ctx, in.ArticleId); err != nil {
			l.Logger.Errorf("IncrViewCount error: %v", err)
		}
		article.ViewCount++
	}

	return &__.GetArticleResponse{
		Article: &__.Article{
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
		},
	}, nil
}
