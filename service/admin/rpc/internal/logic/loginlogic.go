package logic

import (
	"context"
	"errors"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"
	"sea-try-go/service/common/cryptx"

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
	admin := model.Admin{}
	err := l.svcCtx.DB.Where("username = ?", in.Username).First(&admin).Error
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	correct := cryptx.CheckPassword(admin.Password, in.Password)
	if !correct {
		return nil, errors.New("用户名或密码错误")
	}
	return &pb.LoginResp{
		Id: admin.Id,
	}, nil
}
