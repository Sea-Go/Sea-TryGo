package logic

import (
	"context"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/user/rpc/internal/model"
	"sea-try-go/service/user/rpc/internal/svc"
	pb "sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {

	user := model.User{}
	err := l.svcCtx.DB.Where("username = ?", in.Username).First(&user).Error
	if err != nil {
		return &pb.LoginResp{
			Status: 1,
		}, nil
	}

	correct := cryptx.CheckPassword(user.Password, in.Password)
	if !correct {
		return &pb.LoginResp{
			Status: 1,
		}, nil
	}
	if user.Status == 1 {
		return &pb.LoginResp{
			Status: 2,
		}, nil
	}

	return &pb.LoginResp{
		Id:     user.Id,
		Status: 0,
	}, nil
}
