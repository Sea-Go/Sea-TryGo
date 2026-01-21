package logic

import (
	"context"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *pb.GetUserListReq) (*pb.GetUserListResp, error) {
	var users []model.User
	var total int64
	db := l.svcCtx.DB.Model(&model.User{})
	if len(in.Keyword) > 0 {
		keyword := "%" + in.Keyword + "%"
		db = db.Where("username LIKE ? OR email LIKE ?", keyword, keyword)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	list := make([]*pb.UserInfo, 0)
	offset := (in.Page - 1) * in.PageSize
	err := db.Offset(int(offset)).Limit(int(in.PageSize)).Order("id desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		list = append(list, &pb.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Status:    uint64(user.Status),
			ExtraInfo: user.ExtraInfo,
		})
	}
	return &pb.GetUserListResp{
		List:  list,
		Total: total,
	}, nil
}
