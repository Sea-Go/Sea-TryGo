// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"encoding/json"
	"errors"

	"sea-try-go/service/admin/api/internal/model"
	"sea-try-go/service/admin/api/internal/svc"
	"sea-try-go/service/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetselfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetselfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetselfLogic {
	return &GetselfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetselfLogic) Getself(req *types.GetSelfReq) (resp *types.GetSelfResp, err error) {
	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return nil, errors.New("Token 解析异常")
	}
	id, _ := userId.Int64()
	admin := model.Admin{}
	err = l.svcCtx.DB.Where("id = ?", uint64(id)).First(&admin).Error
	return &types.GetSelfResp{
		Admin: types.AdminInfo{
			Id:        admin.Id,
			Username:  admin.Username,
			Email:     admin.Email,
			Extrainfo: admin.ExtraInfo,
		},
	}, nil
}
