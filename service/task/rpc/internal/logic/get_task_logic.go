package logic

import (
	"context"

	"sea-try-go/service/task/rpc/internal/svc"
	"sea-try-go/service/task/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskLogic) GetTask(in *__.GetTaskReq) (*__.GetTaskResp, error) {
	// todo: add your logic here and delete this line

	return &__.GetTaskResp{}, nil
}
