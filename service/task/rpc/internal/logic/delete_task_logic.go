package logic

import (
	"context"

	"sea-try-go/service/task/rpc/internal/svc"
	"sea-try-go/service/task/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogic {
	return &DeleteTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTaskLogic) DeleteTask(in *__.DeleteTaskReq) (*__.DeleteTaskResp, error) {
	// todo: add your logic here and delete this line

	return &__.DeleteTaskResp{}, nil
}
