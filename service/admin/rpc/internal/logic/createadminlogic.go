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

type CreateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAdminLogic) CreateAdmin(in *pb.CreateAdminReq) (*pb.CreateAdminResp, error) {

	var count int64
	l.svcCtx.DB.Model(&model.Admin{}).Where("username = ?", in.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	password, err := cryptx.PasswordEncrypt(in.Password)
	if err != nil {
		return nil, err
	}
	admin := model.Admin{
		Username:  in.Username,
		Password:  password,
		Email:     in.Email,
		ExtraInfo: in.ExtraInfo,
	}
	err = l.svcCtx.DB.Save(&admin).Error
	if err != nil {
		return nil, err
	}
	return &pb.CreateAdminResp{
		Id: admin.Id,
	}, nil
}
