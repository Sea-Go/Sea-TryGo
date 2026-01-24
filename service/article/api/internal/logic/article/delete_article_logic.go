// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"encoding/json"
	"fmt"

	"sea-try-go/service/article/api/internal/svc"
	"sea-try-go/service/article/api/internal/types"
	"sea-try-go/service/article/rpc/articleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticleReq) (resp *types.DeleteArticleResp, err error) {
	var OperatorId string
	if uid := l.ctx.Value("userId"); uid != nil {
		if idNum, ok := uid.(json.Number); ok {
			OperatorId = idNum.String()
		} else {
			OperatorId = fmt.Sprintf("%v", uid)
		}
	} else {
		OperatorId = "dev_test_user"
	}
	_, err = l.svcCtx.ArticleRpc.DeleteArticle(l.ctx, &articleservice.DeleteArticleRequest{
		ArticleId:  req.ArticleId,
		OperatorId: OperatorId,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteArticleResp{
		Success: true,
	}, nil
}
