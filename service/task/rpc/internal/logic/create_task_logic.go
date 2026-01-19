package logic

import (
	"context"

	"sea-try-go/service/task/rpc/internal/svc"
	"sea-try-go/service/task/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTaskLogic) CreateTask(in *__.CreateTaskReq) (*__.CreateTaskResp, error) {
	// todo: add your logic here and delete this line

	return &__.CreateTaskResp{}, nil
}
