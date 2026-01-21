package logic

import (
	"context"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.GetUserResp, error) {
	user := model.User{}
	err := l.svcCtx.DB.Where("id = ?", in.Id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResp{
		User: &pb.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Status:    uint64(user.Status),
			ExtraInfo: user.ExtraInfo,
		},
	}, nil
}
