// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"sea-try-go/service/admin/api/internal/model"
	"sea-try-go/service/admin/api/internal/svc"
	"sea-try-go/service/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetuserlistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetuserlistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetuserlistLogic {
	return &GetuserlistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetuserlistLogic) Getuserlist(req *types.GetUserListReq) (resp *types.GetUserListResp, err error) {
	var users []model.User
	var total int64
	db := l.svcCtx.DB.Model(&model.User{})
	if len(req.Keyword) > 0 {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("username LIKE ? OR email LIKE ?", keyword)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	list := make([]types.UserInfo, 0)
	offset := (req.Page - 1) * req.PageSize
	err = db.Offset(int(offset)).Limit(int(req.PageSize)).Order("id desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		list = append(list, types.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Extrainfo: user.ExtraInfo,
		})
	}
	return &types.GetUserListResp{
		List:  list,
		Total: total,
	}, nil
}
