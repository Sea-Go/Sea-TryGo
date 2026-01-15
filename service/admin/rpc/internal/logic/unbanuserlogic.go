package logic

import (
	"context"
	"errors"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnBanUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnBanUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnBanUserLogic {
	return &UnBanUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnBanUserLogic) UnBanUser(in *pb.UnBanUserReq) (*pb.UnBanUserResp, error) {
	result := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", in.Id).Update("status", 0)
	if result.Error != nil {
		return nil, errors.New("解封失败" + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}
	return &pb.UnBanUserResp{
		Success: true,
	}, nil
}
