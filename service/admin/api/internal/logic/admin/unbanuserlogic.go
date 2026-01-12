// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"

	"sea-try-go/service/admin/api/internal/model"
	"sea-try-go/service/admin/api/internal/svc"
	"sea-try-go/service/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbanuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbanuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbanuserLogic {
	return &UnbanuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbanuserLogic) Unbanuser(req *types.UnBanUserReq) (resp *types.UnBanUserResp, err error) {
	id := req.Id
	result := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", id).Update("status", 0)
	if result.Error != nil {
		return nil, errors.New("解封失败" + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}
	return &types.UnBanUserResp{
		Success: true,
	}, nil
}
