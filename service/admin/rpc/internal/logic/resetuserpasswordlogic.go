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

type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(in *pb.ResetUserPasswordReq) (*pb.ResetUserPasswordResp, error) {

	user := model.User{}
	err := l.svcCtx.DB.Where("id = ?", in.Id).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	var password string
	password, err = cryptx.PasswordEncrypt(l.svcCtx.Config.System.DefaultPassword)
	if err != nil {
		return nil, err
	}
	user.Password = password
	err = l.svcCtx.DB.Model(&user).Update("password", password).Error
	if err != nil {
		return nil, err
	}
	return &pb.ResetUserPasswordResp{
		Success: true,
	}, nil
}
