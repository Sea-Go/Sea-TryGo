package logic

import (
	"context"

	"sea-try-go/service/task/rpc/internal/svc"
	"sea-try-go/service/task/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTaskLogic) UpdateTask(in *__.UpdateTaskReq) (*__.UpdateTaskResp, error) {
	// todo: add your logic here and delete this line

	return &__.UpdateTaskResp{}, nil
}
